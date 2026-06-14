package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"shared/configs/auth"
	"shared/entity"
	core_db "shared/providers/core-db"
)

var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrPasswordTooShort    = errors.New("password too short")
)

const minPasswordLen = 8

// ICartMerger — локальный «view» на cart-usecase: auth знает только то,
// что после логина гостевую корзину нужно слить в пользовательскую.
// Это убирает цикл импорта и фиксирует минимальный контракт зависимости.
type ICartMerger interface {
	MergeGuestIntoUser(ctx context.Context, guestToken string, userID int64) error
}

type IAuthUseCases interface {
	Register(ctx context.Context, in RegisterInUDTO) (*entity.User, *entity.TokenPair, error)
	Login(ctx context.Context, in LoginInUDTO) (*entity.User, *entity.TokenPair, error)
	Refresh(ctx context.Context, in RefreshInUDTO) (*entity.User, *entity.TokenPair, error)
	Logout(ctx context.Context, in LogoutInUDTO) error
	LogoutAll(ctx context.Context, userID int64) error
	RequestPasswordReset(ctx context.Context, in RequestPasswordResetInUDTO) (resetSecret string, _ error)
	ConfirmPasswordReset(ctx context.Context, in ConfirmPasswordResetInUDTO) error
	RequestEmailVerification(ctx context.Context, userID int64) (verifySecret string, _ error)
	ConfirmEmailVerification(ctx context.Context, in ConfirmEmailVerificationInUDTO) (*entity.User, error)
	GetByID(ctx context.Context, userID int64) (*entity.User, error)
}

type useCases struct {
	cfg            *auth.Config
	logger         *zap.Logger
	tx             core_db.ITransactionProvider
	userProv       core_db.IUserProvider
	refreshProv    core_db.IRefreshTokenProvider
	resetProv      core_db.IPasswordResetProvider
	emailVerifyProv core_db.IEmailVerificationProvider
	signer         IAccessTokenSigner
	hasher         IPasswordHasher
	identities     *Registry
	cartMerger     ICartMerger
}

func NewUseCase(
	cfg *auth.Config,
	logger *zap.Logger,
	tx core_db.ITransactionProvider,
	userProv core_db.IUserProvider,
	refreshProv core_db.IRefreshTokenProvider,
	resetProv core_db.IPasswordResetProvider,
	emailVerifyProv core_db.IEmailVerificationProvider,
	signer IAccessTokenSigner,
	hasher IPasswordHasher,
	identities *Registry,
	cartMerger ICartMerger,
) IAuthUseCases {
	return &useCases{
		cfg:             cfg,
		logger:          logger,
		tx:              tx,
		userProv:        userProv,
		refreshProv:     refreshProv,
		resetProv:       resetProv,
		emailVerifyProv: emailVerifyProv,
		signer:          signer,
		hasher:          hasher,
		identities:      identities,
		cartMerger:      cartMerger,
	}
}

func (u *useCases) Register(ctx context.Context, in RegisterInUDTO) (*entity.User, *entity.TokenPair, error) {
	email := normalizeEmail(in.Email)
	if email == "" || !looksLikeEmail(email) {
		return nil, nil, ErrInvalidCredentials
	}
	if len(in.Password) < minPasswordLen {
		return nil, nil, ErrPasswordTooShort
	}

	existing, err := u.userProv.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}
	if existing != nil {
		return nil, nil, ErrEmailAlreadyExists
	}

	hash, err := u.hasher.Hash(in.Password)
	if err != nil {
		return nil, nil, err
	}

	var created *entity.User
	var pair *entity.TokenPair
	err = u.tx.Run(ctx, func(ctx context.Context) error {
		c, err := u.userProv.Create(ctx, entity.User{
			PublicID:     uuid.NewString(),
			Email:        email,
			PasswordHash: &hash,
		})
		if err != nil {
			return err
		}
		created = c

		p, err := u.issueTokenPair(ctx, c, "", "")
		if err != nil {
			return err
		}
		pair = p

		// Сразу выпускаем токен подтверждения email — внутри той же транзакции,
		// чтобы либо юзер+токен оба создались, либо ни один. Сам secret пока
		// никуда не уходит (нет SMTP), но в БД он есть, и при поднятии транспорта
		// добавится отправка письма в одном месте — в issueEmailVerifyToken.
		if _, err := u.issueEmailVerifyToken(ctx, c.ID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return created, pair, nil
}

func (u *useCases) Login(ctx context.Context, in LoginInUDTO) (*entity.User, *entity.TokenPair, error) {
	providerName := in.Provider
	if providerName == "" {
		providerName = ProviderLocal
	}
	provider, err := u.identities.Get(providerName)
	if err != nil {
		return nil, nil, err
	}

	user, err := provider.Authenticate(ctx, Credentials{
		Email:    in.Email,
		Password: in.Password,
	})
	if err != nil {
		return nil, nil, err
	}

	var pair *entity.TokenPair
	err = u.tx.Run(ctx, func(ctx context.Context) error {
		p, err := u.issueTokenPair(ctx, user, in.UserAgent, in.IP)
		if err != nil {
			return err
		}
		pair = p

		if in.GuestToken != "" && u.cartMerger != nil {
			if err := u.cartMerger.MergeGuestIntoUser(ctx, in.GuestToken, user.ID); err != nil {
				// Merge — best-effort; не валим логин, если что-то пошло не так,
				// но логируем — потеря корзины важна для UX.
				u.logger.Warn("guest cart merge failed", zap.Int64("userID", user.ID), zap.Error(err))
			}
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return user, pair, nil
}

func (u *useCases) Refresh(ctx context.Context, in RefreshInUDTO) (*entity.User, *entity.TokenPair, error) {
	if in.RefreshSecret == "" {
		return nil, nil, ErrInvalidRefreshToken
	}
	hash := HashSecret(in.RefreshSecret)

	var pair *entity.TokenPair
	var user *entity.User
	err := u.tx.Run(ctx, func(ctx context.Context) error {
		stored, err := u.refreshProv.GetByHash(ctx, hash)
		if err != nil {
			return err
		}
		if stored == nil || stored.RevokedAt != nil || time.Now().After(stored.ExpiresAt) {
			return ErrInvalidRefreshToken
		}

		u2, err := u.userProv.GetByID(ctx, stored.UserID)
		if err != nil {
			return err
		}
		if u2 == nil {
			return ErrInvalidRefreshToken
		}
		user = u2

		newPair, newID, err := u.issueTokenPairWithID(ctx, user, in.UserAgent, in.IP)
		if err != nil {
			return err
		}
		pair = newPair

		return u.refreshProv.Revoke(ctx, stored.ID, time.Now(), &newID)
	})
	if err != nil {
		return nil, nil, err
	}
	return user, pair, nil
}

func (u *useCases) Logout(ctx context.Context, in LogoutInUDTO) error {
	if in.RefreshSecret == "" {
		return nil
	}
	hash := HashSecret(in.RefreshSecret)
	stored, err := u.refreshProv.GetByHash(ctx, hash)
	if err != nil {
		return err
	}
	if stored == nil || stored.RevokedAt != nil {
		return nil
	}
	return u.refreshProv.Revoke(ctx, stored.ID, time.Now(), nil)
}

func (u *useCases) LogoutAll(ctx context.Context, userID int64) error {
	return u.refreshProv.RevokeAllForUser(ctx, userID, time.Now())
}

func (u *useCases) RequestPasswordReset(ctx context.Context, in RequestPasswordResetInUDTO) (string, error) {
	email := normalizeEmail(in.Email)
	user, err := u.userProv.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		// Не сообщаем "нет такого email" — handler всё равно вернёт 204,
		// чтобы не помогать энумерации.
		return "", nil
	}

	secret, secretHash, err := GenerateRefreshSecret()
	if err != nil {
		return "", err
	}

	ttl := u.cfg.PasswordResetTTL
	if ttl == 0 {
		ttl = time.Hour
	}
	if _, err := u.resetProv.Create(ctx, entity.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: secretHash,
		ExpiresAt: time.Now().Add(ttl),
	}); err != nil {
		return "", err
	}

	// TODO: когда появится email-транспорт, здесь — отправка письма с
	// secret-ссылкой. Сам secret в логи не пишем (это утечка).
	return secret, nil
}

func (u *useCases) ConfirmPasswordReset(ctx context.Context, in ConfirmPasswordResetInUDTO) error {
	if len(in.NewPassword) < minPasswordLen {
		return ErrPasswordTooShort
	}
	hash := HashSecret(in.Secret)

	return u.tx.Run(ctx, func(ctx context.Context) error {
		stored, err := u.resetProv.GetByHash(ctx, hash)
		if err != nil {
			return err
		}
		if stored == nil || stored.UsedAt != nil || time.Now().After(stored.ExpiresAt) {
			return ErrInvalidCredentials
		}

		newHash, err := u.hasher.Hash(in.NewPassword)
		if err != nil {
			return err
		}
		if err := u.userProv.UpdatePasswordHash(ctx, stored.UserID, newHash); err != nil {
			return err
		}
		if err := u.resetProv.MarkUsed(ctx, stored.ID, time.Now()); err != nil {
			return err
		}
		// Безопасный дефолт: после смены пароля все сессии разлогиниваются.
		return u.refreshProv.RevokeAllForUser(ctx, stored.UserID, time.Now())
	})
}

// issueTokenPair выпускает access + refresh, сохраняет refresh-хэш в БД.
// userAgent/ip — для аудита сессий; пустые строки — допустимо (например,
// при регистрации, где смысла привязывать к UA нет).
func (u *useCases) issueTokenPair(ctx context.Context, user *entity.User, userAgent, ip string) (*entity.TokenPair, error) {
	pair, _, err := u.issueTokenPairWithID(ctx, user, userAgent, ip)
	return pair, err
}

func (u *useCases) issueTokenPairWithID(ctx context.Context, user *entity.User, userAgent, ip string) (*entity.TokenPair, int64, error) {
	now := time.Now()

	accessToken, accessExp, err := u.signer.Sign(user.PublicID, user.Email, now)
	if err != nil {
		return nil, 0, err
	}

	refreshSecret, refreshHash, err := GenerateRefreshSecret()
	if err != nil {
		return nil, 0, err
	}

	refreshTTL := u.cfg.RefreshTokenTTL
	if refreshTTL == 0 {
		refreshTTL = 30 * 24 * time.Hour
	}
	refreshExp := now.Add(refreshTTL)

	stored, err := u.refreshProv.Create(ctx, entity.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshHash,
		ExpiresAt: refreshExp,
		UserAgent: nullableString(userAgent),
		IP:        nullableString(ip),
	})
	if err != nil {
		return nil, 0, err
	}

	return &entity.TokenPair{
		AccessToken:        accessToken,
		AccessExpiresAt:    accessExp,
		RefreshTokenSecret: refreshSecret,
		RefreshExpiresAt:   refreshExp,
	}, stored.ID, nil
}

func normalizeEmail(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

func looksLikeEmail(s string) bool {
	at := strings.IndexByte(s, '@')
	return at > 0 && at < len(s)-1 && !strings.ContainsAny(s, " \t\r\n")
}

func nullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// GetByID — для handler'а /auth/me. Возвращает свежие данные пользователя
// (включая email_verified_at) сразу из БД, без побочных эффектов вроде
// rotation refresh-токена.
func (u *useCases) GetByID(ctx context.Context, userID int64) (*entity.User, error) {
	return u.userProv.GetByID(ctx, userID)
}

// RequestEmailVerification — выпускает новый одноразовый токен подтверждения
// email. Старые активные не отзываем намеренно: это «resend», и оба токена
// должны работать до истечения первого (пользователь мог уже нажать ссылку
// в первом письме).
func (u *useCases) RequestEmailVerification(ctx context.Context, userID int64) (string, error) {
	return u.issueEmailVerifyToken(ctx, userID)
}

func (u *useCases) issueEmailVerifyToken(ctx context.Context, userID int64) (string, error) {
	user, err := u.userProv.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrInvalidCredentials
	}
	// Если email уже подтверждён — не плодим токены и возвращаем пустой
	// secret. Handler в этом случае всё равно отдаёт 204 пользователю.
	if user.EmailVerifiedAt != nil {
		return "", nil
	}

	secret, secretHash, err := GenerateRefreshSecret()
	if err != nil {
		return "", err
	}

	ttl := u.cfg.EmailVerifyTTL
	if ttl == 0 {
		ttl = 7 * 24 * time.Hour
	}
	if _, err := u.emailVerifyProv.Create(ctx, entity.EmailVerificationToken{
		UserID:    userID,
		TokenHash: secretHash,
		ExpiresAt: time.Now().Add(ttl),
	}); err != nil {
		return "", err
	}

	// TODO: когда появится email-транспорт — отправка письма со ссылкой
	// /auth/verify-email?token=<secret>. В логи secret НЕ пишем.
	return secret, nil
}

func (u *useCases) ConfirmEmailVerification(ctx context.Context, in ConfirmEmailVerificationInUDTO) (*entity.User, error) {
	if in.Secret == "" {
		return nil, ErrInvalidCredentials
	}
	hash := HashSecret(in.Secret)

	var verified *entity.User
	err := u.tx.Run(ctx, func(ctx context.Context) error {
		stored, err := u.emailVerifyProv.GetByHash(ctx, hash)
		if err != nil {
			return err
		}
		if stored == nil || stored.UsedAt != nil || time.Now().After(stored.ExpiresAt) {
			return ErrInvalidCredentials
		}

		now := time.Now()
		if err := u.userProv.MarkEmailVerified(ctx, stored.UserID, now); err != nil {
			return err
		}
		if err := u.emailVerifyProv.MarkUsed(ctx, stored.ID, now); err != nil {
			return err
		}

		u2, err := u.userProv.GetByID(ctx, stored.UserID)
		if err != nil {
			return err
		}
		verified = u2
		return nil
	})
	if err != nil {
		return nil, err
	}
	return verified, nil
}