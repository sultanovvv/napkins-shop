-- Корзина. Одна строка одновременно либо пользовательская (user_id),
-- либо гостевая (guest_token_hash от opaque-секрета в HttpOnly-cookie).
-- CHECK гарантирует "ровно один" из двух; частичные UNIQUE-индексы дают
-- "у пользователя/гостя — одна корзина".
CREATE TABLE IF NOT EXISTS carts
(
    id                BIGSERIAL PRIMARY KEY,
    user_id           BIGINT REFERENCES users (id) ON DELETE CASCADE,
    guest_token_hash  TEXT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK ((user_id IS NULL) <> (guest_token_hash IS NULL))
);

CREATE UNIQUE INDEX IF NOT EXISTS carts_user_uniq_idx
    ON carts (user_id) WHERE user_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS carts_guest_uniq_idx
    ON carts (guest_token_hash) WHERE guest_token_hash IS NOT NULL;

-- Позиции в корзине. PK по (cart_id, product_id) автоматически даёт
-- "не больше одной строки на товар"; quantity обновляется суммированием.
CREATE TABLE IF NOT EXISTS cart_items
(
    cart_id    BIGINT      NOT NULL REFERENCES carts (id) ON DELETE CASCADE,
    product_id BIGINT      NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    quantity   INT         NOT NULL CHECK (quantity > 0),
    added_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (cart_id, product_id)
);

CREATE INDEX IF NOT EXISTS cart_items_product_idx ON cart_items (product_id);