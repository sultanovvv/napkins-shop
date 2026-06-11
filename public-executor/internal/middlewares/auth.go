package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	authuc "napkins-shop/public-executor/internal/usecases/auth"
	core_db "shared/providers/core-db"
)

// Ключи в echo.Context. Унифицированы, чтобы handler'ы могли вытащить
// текущего юзера/гостя без знания о JWT и cookie.
const (
	CtxUserID       = "auth.user_id"
	CtxUserPublicID = "auth.user_public_id"
	CtxUserEmail    = "auth.user_email"
	CtxGuestToken   = "auth.guest_token"
)

// AuthMiddleware — фабрика для access-JWT-middleware. Конструируется один раз
// в DI; даёт два варианта применения на route.
type AuthMiddleware struct {
	signer   authuc.IAccessTokenSigner
	userProv core_db.IUserProvider
}

func NewAuthMiddleware(signer authuc.IAccessTokenSigner, userProv core_db.IUserProvider) *AuthMiddleware {
	return &AuthMiddleware{signer: signer, userProv: userProv}
}

// RequireUser — 401, если токена нет или он невалиден. Кладёт user_id в ctx.
func (m *AuthMiddleware) RequireUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ok, err := m.populate(c)
			if err != nil {
				return unauthorized(c, err.Error())
			}
			if !ok {
				return unauthorized(c, "authentication required")
			}
			return next(c)
		}
	}
}

// OptionalUser — если токен есть и валиден, кладёт user_id; иначе тихо
// пропускает. Нужен на /cart/* — туда ходят и юзеры, и гости.
func (m *AuthMiddleware) OptionalUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, _ = m.populate(c)
			return next(c)
		}
	}
}

func (m *AuthMiddleware) populate(c echo.Context) (bool, error) {
	raw := bearerFromHeader(c.Request().Header.Get("Authorization"))
	if raw == "" {
		return false, nil
	}
	claims, err := m.signer.Parse(raw)
	if err != nil {
		if errors.Is(err, authuc.ErrExpiredAccessToken) {
			return false, errors.New("access token expired")
		}
		return false, errors.New("invalid access token")
	}

	user, err := m.userProv.GetByPublicID(c.Request().Context(), claims.Subject)
	if err != nil || user == nil {
		return false, errors.New("invalid access token")
	}

	c.Set(CtxUserID, user.ID)
	c.Set(CtxUserPublicID, user.PublicID)
	c.Set(CtxUserEmail, user.Email)
	// Прокидываем user_id в request.Context — на случай если usecase будет
	// доставать его без явной передачи (для логов, аудита и т.п.).
	ctx := context.WithValue(c.Request().Context(), ctxKeyUserID{}, user.ID)
	c.SetRequest(c.Request().WithContext(ctx))
	return true, nil
}

func bearerFromHeader(h string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(h, prefix) {
		return ""
	}
	return strings.TrimSpace(h[len(prefix):])
}

func unauthorized(c echo.Context, msg string) error {
	return c.JSON(http.StatusUnauthorized, map[string]any{
		"error": map[string]string{
			"code":    "UNAUTHORIZED",
			"message": msg,
		},
	})
}

type ctxKeyUserID struct{}

// UserIDFromContext — удобный доступ из usecase'ов, если когда-нибудь
// понадобится передавать "текущего юзера" через context, а не аргументом.
func UserIDFromContext(ctx context.Context) (int64, bool) {
	v, ok := ctx.Value(ctxKeyUserID{}).(int64)
	return v, ok
}
