package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"shared/configs/auth"
)

// AccessClaims — payload JWT. sub = user.public_id, чтобы наружу никогда не
// светить внутренний BIGSERIAL.
type AccessClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

var (
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrExpiredAccessToken = errors.New("expired access token")
)

type IAccessTokenSigner interface {
	Sign(userPublicID, email string, now time.Time) (token string, expiresAt time.Time, err error)
	Parse(token string) (*AccessClaims, error)
}

type hs256Signer struct {
	secret []byte
	ttl    time.Duration
}

func NewHS256TokenSigner(cfg *auth.Config) IAccessTokenSigner {
	ttl := cfg.AccessTokenTTL
	if ttl == 0 {
		ttl = 15 * time.Minute
	}
	return &hs256Signer{
		secret: []byte(cfg.JWTSecret),
		ttl:    ttl,
	}
}

func (s *hs256Signer) Sign(userPublicID, email string, now time.Time) (string, time.Time, error) {
	expiresAt := now.Add(s.ttl)
	claims := AccessClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userPublicID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign jwt: %w", err)
	}
	return signed, expiresAt, nil
}

func (s *hs256Signer) Parse(token string) (*AccessClaims, error) {
	claims := &AccessClaims{}
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidAccessToken
		}
		return s.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredAccessToken
		}
		return nil, ErrInvalidAccessToken
	}
	if !t.Valid {
		return nil, ErrInvalidAccessToken
	}
	return claims, nil
}

// refreshSecretBytes — 32 байта = 256 бит энтропии. URL-safe base64 без
// паддинга, чтобы спокойно класть в cookie/header.
const refreshSecretBytes = 32

// GenerateRefreshSecret возвращает (secret, secretHash). Secret уходит клиенту
// (cookie), secretHash — в БД. Из БД восстановить secret невозможно.
func GenerateRefreshSecret() (secret, secretHash string, err error) {
	buf := make([]byte, refreshSecretBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", "", fmt.Errorf("read random: %w", err)
	}
	secret = base64.RawURLEncoding.EncodeToString(buf)
	secretHash = HashSecret(secret)
	return secret, secretHash, nil
}

// HashSecret — SHA-256 в hex. Используется и для refresh-токенов, и для
// password-reset, и для guest-cart-token: один и тот же класс "opaque-секрет
// в cookie/email, хэш в БД".
func HashSecret(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:])
}
