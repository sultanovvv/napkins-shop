-- Refresh-токены. В БД хранится ТОЛЬКО SHA-256 хэш от секрета (32 байта
-- crypto/rand → 256 бит энтропии), сам секрет уходит клиенту в HttpOnly-cookie
-- и в БД больше никогда не появляется. replaced_by_id даёт аудит-цепочку
-- rotation: при /refresh старый токен помечается revoked_at и связывается
-- с новым через replaced_by_id.
CREATE TABLE IF NOT EXISTS refresh_tokens
(
    id             BIGSERIAL PRIMARY KEY,
    user_id        BIGINT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash     TEXT        NOT NULL UNIQUE,
    issued_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at     TIMESTAMPTZ NOT NULL,
    revoked_at     TIMESTAMPTZ,
    replaced_by_id BIGINT REFERENCES refresh_tokens (id) ON DELETE SET NULL,
    user_agent     TEXT,
    ip             TEXT
);

-- Активные токены пользователя ищутся постоянно (logout-all, проверка лимита
-- сессий), частичный индекс заметно быстрее полного.
CREATE INDEX IF NOT EXISTS refresh_tokens_active_user_idx
    ON refresh_tokens (user_id) WHERE revoked_at IS NULL;

-- Для будущего cleanup-job (вычистить просроченные неотозванные).
CREATE INDEX IF NOT EXISTS refresh_tokens_active_expiry_idx
    ON refresh_tokens (expires_at) WHERE revoked_at IS NULL;

-- Токены восстановления пароля. Одноразовые, короткоживущие. Хранится только
-- хэш. used_at защищает от повторного применения (race + replay).
CREATE TABLE IF NOT EXISTS password_reset_tokens
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash TEXT        NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at    TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS password_reset_tokens_user_idx ON password_reset_tokens (user_id);