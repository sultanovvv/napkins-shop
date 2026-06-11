package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"

	"shared/configs/auth"
)

// PHC-string формата $argon2id$v=19$m=...,t=...,p=...$<salt-b64>$<hash-b64>.
// Параметры лежат внутри строки, поэтому при изменении настроек старые хэши
// продолжают валидироваться по своим параметрам, новые — пишутся по новым.

var (
	ErrInvalidPasswordHash = errors.New("invalid password hash format")
	ErrUnsupportedHashAlgo = errors.New("unsupported password hash algorithm")
)

type IPasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, encodedHash string) (bool, error)
}

type argon2Hasher struct {
	cfg auth.Argon2Config
}

func NewArgon2Hasher(cfg auth.Argon2Config) IPasswordHasher {
	if cfg.Memory == 0 {
		cfg.Memory = 64 * 1024
	}
	if cfg.Iterations == 0 {
		cfg.Iterations = 3
	}
	if cfg.Parallelism == 0 {
		cfg.Parallelism = 2
	}
	if cfg.SaltLen == 0 {
		cfg.SaltLen = 16
	}
	if cfg.KeyLen == 0 {
		cfg.KeyLen = 32
	}
	return &argon2Hasher{cfg: cfg}
}

func (h *argon2Hasher) Hash(password string) (string, error) {
	salt := make([]byte, h.cfg.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("read salt: %w", err)
	}
	key := argon2.IDKey([]byte(password), salt, h.cfg.Iterations, h.cfg.Memory, h.cfg.Parallelism, h.cfg.KeyLen)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		h.cfg.Memory, h.cfg.Iterations, h.cfg.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

func (h *argon2Hasher) Verify(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[0] != "" {
		return false, ErrInvalidPasswordHash
	}
	if parts[1] != "argon2id" {
		return false, ErrUnsupportedHashAlgo
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return false, ErrInvalidPasswordHash
	}
	if version != argon2.Version {
		return false, ErrUnsupportedHashAlgo
	}

	var memory, iterations uint32
	var parallelism uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism); err != nil {
		return false, ErrInvalidPasswordHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, ErrInvalidPasswordHash
	}
	expected, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, ErrInvalidPasswordHash
	}

	got := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(expected)))
	return subtle.ConstantTimeCompare(got, expected) == 1, nil
}