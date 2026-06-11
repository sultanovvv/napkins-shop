package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"napkins-shop/public-executor/api"
	"napkins-shop/public-executor/internal/middlewares"
	authuc "napkins-shop/public-executor/internal/usecases/auth"
	authcfg "shared/configs/auth"
	"shared/entity"
)

// refreshCookieName и refreshCookiePath — refresh-секрет ходит только на
// /api/v1/public/auth/refresh, чтобы он не отправлялся вместе с обычными
// запросами и не уходил в чужие хендлеры.
const (
	refreshCookieName = "rft"
	refreshCookiePath = "/api/v1/public/auth/refresh"
)

type AuthHandler struct {
	useCase authuc.IAuthUseCases
	logger  *zap.Logger
	cfg     *authcfg.Config
	secure  bool
}

func NewAuthHandler(useCase authuc.IAuthUseCases, cfg *authcfg.Config, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		useCase: useCase,
		cfg:     cfg,
		logger:  logger,
		secure:  false, // dev-дефолт; в проде переключить через конфиг/прокси
	}
}

func (h *AuthHandler) AuthRegister(ctx echo.Context) error {
	var body api.AuthRegisterRequest
	if err := ctx.Bind(&body); err != nil {
		return badRequest(ctx, "invalid request body")
	}

	user, pair, err := h.useCase.Register(ctx.Request().Context(), authuc.RegisterInUDTO{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, authuc.ErrEmailAlreadyExists):
			return errorJSON(ctx, http.StatusConflict, "EMAIL_EXISTS", "email already registered")
		case errors.Is(err, authuc.ErrPasswordTooShort):
			return badRequest(ctx, "password too short")
		case errors.Is(err, authuc.ErrInvalidCredentials):
			return badRequest(ctx, "invalid email or password")
		default:
			h.logger.Error("register failed", zap.Error(err))
			return internalErr(ctx)
		}
	}

	writeRefreshCookie(ctx, pair, h.secure)
	return ctx.JSON(http.StatusOK, loginResponseFromEntity(user, pair))
}

func (h *AuthHandler) AuthLogin(ctx echo.Context) error {
	var body api.AuthLoginRequest
	if err := ctx.Bind(&body); err != nil {
		return badRequest(ctx, "invalid request body")
	}

	in := authuc.LoginInUDTO{
		Email:      body.Email,
		Password:   body.Password,
		GuestToken: readGuestTokenFromCtx(ctx),
		UserAgent:  ctx.Request().UserAgent(),
		IP:         ctx.RealIP(),
	}
	if body.Provider != nil {
		in.Provider = *body.Provider
	}

	user, pair, err := h.useCase.Login(ctx.Request().Context(), in)
	if err != nil {
		switch {
		case errors.Is(err, authuc.ErrInvalidCredentials),
			errors.Is(err, authuc.ErrPasswordLoginDisabled),
			errors.Is(err, authuc.ErrUnknownIdentityProvider):
			return unauthorizedJSON(ctx, "invalid credentials")
		default:
			h.logger.Error("login failed", zap.Error(err))
			return internalErr(ctx)
		}
	}

	// Гостевая корзина сметана в пользовательскую — обнуляем gct-cookie,
	// чтобы при будущей разлогинке не цеплять призрак удалённой корзины.
	if in.GuestToken != "" {
		middlewares.ClearGuestCookie(ctx)
	}

	writeRefreshCookie(ctx, pair, h.secure)
	return ctx.JSON(http.StatusOK, loginResponseFromEntity(user, pair))
}

func (h *AuthHandler) AuthRefresh(ctx echo.Context) error {
	secret := readRefreshCookie(ctx)
	if secret == "" {
		return unauthorizedJSON(ctx, "missing refresh token")
	}

	user, pair, err := h.useCase.Refresh(ctx.Request().Context(), authuc.RefreshInUDTO{
		RefreshSecret: secret,
		UserAgent:     ctx.Request().UserAgent(),
		IP:            ctx.RealIP(),
	})
	if err != nil {
		if errors.Is(err, authuc.ErrInvalidRefreshToken) {
			clearRefreshCookie(ctx)
			return unauthorizedJSON(ctx, "invalid refresh token")
		}
		h.logger.Error("refresh failed", zap.Error(err))
		return internalErr(ctx)
	}

	writeRefreshCookie(ctx, pair, h.secure)
	return ctx.JSON(http.StatusOK, loginResponseFromEntity(user, pair))
}

func (h *AuthHandler) AuthLogout(ctx echo.Context) error {
	secret := readRefreshCookie(ctx)
	if err := h.useCase.Logout(ctx.Request().Context(), authuc.LogoutInUDTO{
		RefreshSecret: secret,
	}); err != nil {
		h.logger.Error("logout failed", zap.Error(err))
		return internalErr(ctx)
	}
	clearRefreshCookie(ctx)
	return ctx.NoContent(http.StatusNoContent)
}

func (h *AuthHandler) AuthLogoutAll(ctx echo.Context) error {
	userID, ok := ctx.Get(middlewares.CtxUserID).(int64)
	if !ok {
		return unauthorizedJSON(ctx, "authentication required")
	}
	if err := h.useCase.LogoutAll(ctx.Request().Context(), userID); err != nil {
		h.logger.Error("logout-all failed", zap.Error(err))
		return internalErr(ctx)
	}
	clearRefreshCookie(ctx)
	return ctx.NoContent(http.StatusNoContent)
}

func (h *AuthHandler) AuthRequestPasswordReset(ctx echo.Context) error {
	var body api.AuthRequestPasswordResetRequest
	if err := ctx.Bind(&body); err != nil {
		return badRequest(ctx, "invalid request body")
	}
	if _, err := h.useCase.RequestPasswordReset(ctx.Request().Context(), authuc.RequestPasswordResetInUDTO{
		Email: body.Email,
	}); err != nil {
		// Логируем, но клиенту всё равно 204 — не помогаем энумерации.
		h.logger.Error("request password reset failed", zap.Error(err))
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (h *AuthHandler) AuthMe(ctx echo.Context) error {
	uid, ok := ctx.Get(middlewares.CtxUserID).(int64)
	if !ok {
		return unauthorizedJSON(ctx, "authentication required")
	}
	user, err := h.useCase.GetByID(ctx.Request().Context(), uid)
	if err != nil {
		h.logger.Error("me failed", zap.Error(err))
		return internalErr(ctx)
	}
	if user == nil {
		return unauthorizedJSON(ctx, "user not found")
	}
	return ctx.JSON(http.StatusOK, authUserResponse(user))
}

func (h *AuthHandler) AuthRequestEmailVerification(ctx echo.Context) error {
	uid, ok := ctx.Get(middlewares.CtxUserID).(int64)
	if !ok {
		return unauthorizedJSON(ctx, "authentication required")
	}
	if _, err := h.useCase.RequestEmailVerification(ctx.Request().Context(), uid); err != nil {
		// Логируем, но клиенту 204 — поведение симметрично RequestPasswordReset.
		h.logger.Error("request email verification failed", zap.Error(err))
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (h *AuthHandler) AuthConfirmEmailVerification(ctx echo.Context) error {
	var body api.AuthConfirmEmailVerificationRequest
	if err := ctx.Bind(&body); err != nil {
		return badRequest(ctx, "invalid request body")
	}
	user, err := h.useCase.ConfirmEmailVerification(ctx.Request().Context(), authuc.ConfirmEmailVerificationInUDTO{
		Secret: body.Secret,
	})
	if err != nil {
		if errors.Is(err, authuc.ErrInvalidCredentials) {
			return badRequest(ctx, "invalid or expired verification token")
		}
		h.logger.Error("confirm email verification failed", zap.Error(err))
		return internalErr(ctx)
	}
	return ctx.JSON(http.StatusOK, authUserResponse(user))
}

func (h *AuthHandler) AuthConfirmPasswordReset(ctx echo.Context) error {
	var body api.AuthConfirmPasswordResetRequest
	if err := ctx.Bind(&body); err != nil {
		return badRequest(ctx, "invalid request body")
	}
	err := h.useCase.ConfirmPasswordReset(ctx.Request().Context(), authuc.ConfirmPasswordResetInUDTO{
		Secret:      body.Secret,
		NewPassword: body.NewPassword,
	})
	if err != nil {
		switch {
		case errors.Is(err, authuc.ErrInvalidCredentials):
			return badRequest(ctx, "invalid or expired reset token")
		case errors.Is(err, authuc.ErrPasswordTooShort):
			return badRequest(ctx, "password too short")
		default:
			h.logger.Error("confirm password reset failed", zap.Error(err))
			return internalErr(ctx)
		}
	}
	return ctx.NoContent(http.StatusNoContent)
}

func writeRefreshCookie(ctx echo.Context, pair *entity.TokenPair, secure bool) {
	ctx.SetCookie(&http.Cookie{
		Name:     refreshCookieName,
		Value:    pair.RefreshTokenSecret,
		Path:     refreshCookiePath,
		Expires:  pair.RefreshExpiresAt,
		MaxAge:   int(time.Until(pair.RefreshExpiresAt).Seconds()),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearRefreshCookie(ctx echo.Context) {
	ctx.SetCookie(&http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     refreshCookiePath,
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func readRefreshCookie(ctx echo.Context) string {
	ck, err := ctx.Cookie(refreshCookieName)
	if err != nil || ck == nil {
		return ""
	}
	return ck.Value
}

func readGuestTokenFromCtx(ctx echo.Context) string {
	v, _ := ctx.Get(middlewares.CtxGuestToken).(string)
	return v
}