BEGIN;

ALTER TABLE products
    ADD COLUMN description TEXT NOT NULL DEFAULT '',
    ADD COLUMN price_minor BIGINT,
    ADD COLUMN currency_code VARCHAR(3),
    ADD COLUMN stock_quantity INTEGER,
    ADD COLUMN is_published BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    ADD CONSTRAINT products_name_not_blank CHECK (length(btrim(name)) > 0),
    ADD CONSTRAINT products_slug_not_blank CHECK (length(btrim(slug)) > 0),
    ADD CONSTRAINT products_price_nonnegative CHECK (price_minor >= 0),
    ADD CONSTRAINT products_stock_nonnegative CHECK (stock_quantity >= 0),
    ADD CONSTRAINT products_price_currency_pair
        CHECK ((price_minor IS NULL) = (currency_code IS NULL)),
    ADD CONSTRAINT products_currency_format CHECK (currency_code ~ '^[A-Z]{3}$'),
    ADD CONSTRAINT products_publication_complete CHECK (
        NOT is_published OR (
            price_minor IS NOT NULL
            AND currency_code IS NOT NULL
            AND stock_quantity IS NOT NULL
        )
    );

-- Старые цены и остатки неизвестны: NULL сохраняет это отличие от нуля.
-- Существующие товары остаются в БД как неопубликованные черновики.
COMMENT ON COLUMN products.price_minor IS 'Цена в минимальных денежных единицах; NULL — не заполнена';
COMMENT ON COLUMN products.stock_quantity IS 'Остаток в целых штуках; NULL — не заполнен';
COMMENT ON COLUMN products.created_at IS 'Для существовавших до миграции товаров — дата переноса';

COMMIT;
