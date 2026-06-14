package auth

import (
	"context"
	"errors"
	"strings"

	"shared/entity"
	core_db "shared/providers/core-db"
)

// IIdentityProvider — абстракция над способом подтверждения, что у запроса
// есть владелец. У каждой реализации свой набор обязательных полей в
// Credentials (см. соответствующий провайдер).
//
// Возвращает entity.User: если пользователь ещё не зарегистрирован в нашей
// БД (актуально только для внешних провайдеров), реализация сама создаёт
// запись в users + auth_identities. Local-провайдер юзера не создаёт —
// для него регистрация — отдельный явный шаг (AuthUseCases.Register).
type IIdentityProvider interface {
	Name() string
	Authenticate(ctx context.Context, c Credentials) (*entity.User, error)
}

// Credentials — типобезопасных вариантов на каждый провайдер не делаем,
// чтобы регистр был плоский. Внешние провайдеры читают свои поля; неизвестные
// игнорируют.
type Credentials struct {
	Email    string // local
	Password string // local
	IDToken  string // future: OIDC providers
	Code     string // future: OAuth2 code-flow
}

var (
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrPasswordLoginDisabled   = errors.New("password login disabled for this user")
	ErrUnknownIdentityProvider = errors.New("unknown identity provider")
)

// Registry — реестр доступных IdentityProvider'ов по имени. Local
// регистрируется на старте. Google/Yandex/VK добавятся когда появятся
// реализации — без изменений в usecase/handler'ах.
type Registry struct {
	providers map[string]IIdentityProvider
}

func NewRegistry(providers ...IIdentityProvider) *Registry {
	m := make(map[string]IIdentityProvider, len(providers))
	for _, p := range providers {
		m[p.Name()] = p
	}
	return &Registry{providers: m}
}

func (r *Registry) Get(name string) (IIdentityProvider, error) {
	p, ok := r.providers[name]
	if !ok {
		return nil, ErrUnknownIdentityProvider
	}
	return p, nil
}

// ProviderLocal — имя локального (email/password) провайдера. Константа,
// чтобы не было опечаток в регистрах.
const ProviderLocal = "local"

type localIdentityProvider struct {
	userProv core_db.IUserProvider
	hasher   IPasswordHasher
}

func NewLocalIdentityProvider(userProv core_db.IUserProvider, hasher IPasswordHasher) IIdentityProvider {
	return &localIdentityProvider{userProv: userProv, hasher: hasher}
}

func (p *localIdentityProvider) Name() string { return ProviderLocal }

func (p *localIdentityProvider) Authenticate(ctx context.Context, c Credentials) (*entity.User, error) {
	email := strings.TrimSpace(c.Email)
	if email == "" || c.Password == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := p.userProv.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		// Не сообщаем "пользователь не найден" — единое сообщение для логина,
		// чтобы не помогать энумерации email'ов.
		return nil, ErrInvalidCredentials
	}
	if user.PasswordHash == nil {
		return nil, ErrPasswordLoginDisabled
	}

	ok, err := p.hasher.Verify(c.Password, *user.PasswordHash)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}
