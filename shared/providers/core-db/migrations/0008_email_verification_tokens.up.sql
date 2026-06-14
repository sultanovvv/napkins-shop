-- Токены подтверждения email. Структура близка к password_reset_tokens:
-- хранится только хэш, secret уходит в письмо (когда появится транспорт).
-- Подтверждение записывает users.email_verified_at.
CREATE TABLE IF NOT EXISTS email_verification_tokens
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash TEXT        NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at    TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS email_verification_tokens_user_idx ON email_verification_tokens (user_id);
