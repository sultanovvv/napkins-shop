-- Пользователи. id внутренний, public_id (UUID) — публичный идентификатор
-- (попадает в sub JWT, не утечёт порядок регистраций). Email хранится как
-- введён, но уникальность — по LOWER(email), чтобы Vasya@x.ru и vasya@x.ru
-- считались одним юзером.
--
-- password_hash NULLABLE: при чистой регистрации через будущий OIDC-провайдер
-- локального пароля у пользователя может не быть. Тогда логин по паролю
-- невозможен, а внешний логин идёт через auth_identities.
CREATE TABLE IF NOT EXISTS users
(
    id                BIGSERIAL PRIMARY KEY,
    public_id         UUID        NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    email             TEXT        NOT NULL,
    password_hash     TEXT,
    email_verified_at TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS users_email_lower_idx ON users (LOWER(email));

-- Внешние identity-провайдеры (Google/Yandex/VK/...). Запись на пару
-- (provider, provider_user_id). Локальный пароль здесь НЕ хранится — для
-- local-логина признак "у юзера есть пароль" = users.password_hash IS NOT NULL.
-- Это убирает дубликат email ↔ provider_user_id и упрощает миграции.
CREATE TABLE IF NOT EXISTS auth_identities
(
    id               BIGSERIAL PRIMARY KEY,
    user_id          BIGINT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    provider         TEXT        NOT NULL,
    provider_user_id TEXT        NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (provider, provider_user_id)
);

CREATE INDEX IF NOT EXISTS auth_identities_user_idx ON auth_identities (user_id);