CREATE TYPE attribute_value_type AS ENUM ('string', 'int');

-- Каталог доступных атрибутов товара. Закрытый набор: расширение —
-- новая строка через миграцию. value_type подсказывает приложению, какое
-- поле читать из product_attribute_values (text vs int).
CREATE TABLE IF NOT EXISTS attributes
(
    id         BIGSERIAL PRIMARY KEY,
    slug       TEXT                 NOT NULL UNIQUE,
    name       TEXT                 NOT NULL,
    value_type attribute_value_type NOT NULL,
    sort_order INT                  NOT NULL DEFAULT 0
);

INSERT INTO attributes (slug, name, value_type, sort_order)
VALUES ('size', 'Размер', 'string', 10),
       ('layers', 'Слои', 'int', 20),
       ('color', 'Цвет', 'string', 30),
       ('theme', 'Тема', 'string', 40),
       ('pattern', 'Узор', 'string', 50)
ON CONFLICT (slug) DO NOTHING;

-- Значения атрибутов для конкретных товаров. Композитный PK даёт
-- "не более одного значения каждого атрибута на товар".
CREATE TABLE IF NOT EXISTS product_attribute_values
(
    product_id   BIGINT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    attribute_id BIGINT NOT NULL REFERENCES attributes (id) ON DELETE CASCADE,
    value_text   TEXT,
    value_int    INT,
    PRIMARY KEY (product_id, attribute_id),
    CHECK ((value_text IS NULL) <> (value_int IS NULL))
);

CREATE INDEX IF NOT EXISTS pav_attribute_idx ON product_attribute_values (attribute_id);

INSERT INTO product_attribute_values (product_id, attribute_id, value_text, value_int)
VALUES
    -- roses-napkin
    ((SELECT id FROM products WHERE slug = 'roses-napkin'),
     (SELECT id FROM attributes WHERE slug = 'size'), '33x33 см', NULL),
    ((SELECT id FROM products WHERE slug = 'roses-napkin'),
     (SELECT id FROM attributes WHERE slug = 'layers'), NULL, 3),
    ((SELECT id FROM products WHERE slug = 'roses-napkin'),
     (SELECT id FROM attributes WHERE slug = 'color'), 'белый', NULL),
    ((SELECT id FROM products WHERE slug = 'roses-napkin'),
     (SELECT id FROM attributes WHERE slug = 'pattern'), 'розы', NULL),
    -- daisies-napkin
    ((SELECT id FROM products WHERE slug = 'daisies-napkin'),
     (SELECT id FROM attributes WHERE slug = 'size'), '25x25 см', NULL),
    ((SELECT id FROM products WHERE slug = 'daisies-napkin'),
     (SELECT id FROM attributes WHERE slug = 'layers'), NULL, 2),
    ((SELECT id FROM products WHERE slug = 'daisies-napkin'),
     (SELECT id FROM attributes WHERE slug = 'pattern'), 'цветочный узор', NULL),
    -- butterflies-napkin
    ((SELECT id FROM products WHERE slug = 'butterflies-napkin'),
     (SELECT id FROM attributes WHERE slug = 'size'), '33x33 см', NULL),
    ((SELECT id FROM products WHERE slug = 'butterflies-napkin'),
     (SELECT id FROM attributes WHERE slug = 'layers'), NULL, 3),
    ((SELECT id FROM products WHERE slug = 'butterflies-napkin'),
     (SELECT id FROM attributes WHERE slug = 'color'), 'голубой', NULL),
    ((SELECT id FROM products WHERE slug = 'butterflies-napkin'),
     (SELECT id FROM attributes WHERE slug = 'pattern'), 'бабочки', NULL),
    -- leaves-napkin
    ((SELECT id FROM products WHERE slug = 'leaves-napkin'),
     (SELECT id FROM attributes WHERE slug = 'size'), '25x25 см', NULL),
    ((SELECT id FROM products WHERE slug = 'leaves-napkin'),
     (SELECT id FROM attributes WHERE slug = 'layers'), NULL, 2),
    ((SELECT id FROM products WHERE slug = 'leaves-napkin'),
     (SELECT id FROM attributes WHERE slug = 'color'), 'зелёный', NULL),
    ((SELECT id FROM products WHERE slug = 'leaves-napkin'),
     (SELECT id FROM attributes WHERE slug = 'theme'), 'осенний', NULL),
    ((SELECT id FROM products WHERE slug = 'leaves-napkin'),
     (SELECT id FROM attributes WHERE slug = 'pattern'), 'листья', NULL),
    -- gnomes-napkin
    ((SELECT id FROM products WHERE slug = 'gnomes-napkin'),
     (SELECT id FROM attributes WHERE slug = 'size'), '33x33 см', NULL),
    ((SELECT id FROM products WHERE slug = 'gnomes-napkin'),
     (SELECT id FROM attributes WHERE slug = 'layers'), NULL, 3),
    ((SELECT id FROM products WHERE slug = 'gnomes-napkin'),
     (SELECT id FROM attributes WHERE slug = 'theme'), 'новогодний', NULL),
    ((SELECT id FROM products WHERE slug = 'gnomes-napkin'),
     (SELECT id FROM attributes WHERE slug = 'pattern'), 'гномы', NULL)
ON CONFLICT (product_id, attribute_id) DO NOTHING;
