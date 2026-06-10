-- Метаданные изображений товара. Бинарь живёт в S3 (Yandex Object Storage),
-- здесь только s3_key и порядок. URL собирает приложение из S3-конфига.
CREATE TABLE IF NOT EXISTS product_images
(
    id         BIGSERIAL PRIMARY KEY,
    product_id BIGINT      NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    s3_key     TEXT        NOT NULL,
    sort_order INT         NOT NULL DEFAULT 0,
    is_primary BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (product_id, s3_key)
);

CREATE INDEX IF NOT EXISTS product_images_product_idx ON product_images (product_id, sort_order);

-- Не больше одной primary-картинки на товар.
CREATE UNIQUE INDEX IF NOT EXISTS product_images_one_primary_idx
    ON product_images (product_id) WHERE is_primary;
