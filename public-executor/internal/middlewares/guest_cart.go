package middlewares

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	authuc "napkins-shop/public-executor/internal/usecases/auth"
	"shared/configs/auth"
)

// GuestCartCookieName — opaque-секрет, серверу важен только его SHA-256.
const GuestCartCookieName = "gct"

// GuestCartMiddleware читает (а при отсутствии — выпускает) гостевой токен
// в HttpOnly-cookie. Кладёт сырой токен в echo-context (CtxGuestToken),
// чтобы handler/usecase могли посчитать SHA-256 и найти/создать корзину.
//
// Cookie выдаётся ВСЕМ, даже залогиненным: упрощает мердж при следующей
// разлогинке. Авторизованный путь в usecase сам игнорирует guest-token, если
// есть user_id.
type GuestCartMiddleware struct {
	ttl    time.Duration
	secure bool
}

func NewGuestCartMiddleware(cfg *auth.Config) *GuestCartMiddleware {
	ttl := cfg.GuestCartTTL
	if ttl == 0 {
		ttl = 180 * 24 * time.Hour
	}
	return &GuestCartMiddleware{
		ttl: ttl,
		// Secure=false для локальной разработки на http://localhost. На проде
		// прокси проставит "X-Forwarded-Proto: https", и можно будет включить
		// — но сейчас фиксируем минимально совместимый дефолт.
		secure: false,
	}
}

func (m *GuestCartMiddleware) Handle() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := readGuestCookie(c)
			if token == "" {
				newSecret, _, err := authuc.GenerateRefreshSecret()
				if err != nil {
					return next(c)
				}
				token = newSecret
				setGuestCookie(c, token, m.ttl, m.secure)
			}
			c.Set(CtxGuestToken, token)
			return next(c)
		}
	}
}

func readGuestCookie(c echo.Context) string {
	ck, err := c.Cookie(GuestCartCookieName)
	if err != nil || ck == nil {
		return ""
	}
	return ck.Value
}

func setGuestCookie(c echo.Context, value string, ttl time.Duration, secure bool) {
	c.SetCookie(&http.Cookie{
		Name:     GuestCartCookieName,
		Value:    value,
		Path:     "/",
		MaxAge:   int(ttl.Seconds()),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearGuestCookie — handler /auth/login может выкинуть guest-cookie после
// успешного merge, чтобы повторный логин не цеплял удалённую гостевую
// корзину.
func ClearGuestCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     GuestCartCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
