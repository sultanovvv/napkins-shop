CREATE TABLE IF NOT EXISTS categories
(
    id         BIGSERIAL PRIMARY KEY,
    slug       TEXT   NOT NULL UNIQUE,
    name       TEXT   NOT NULL,
    parent_id  BIGINT REFERENCES categories (id) ON DELETE CASCADE,
    sort_order INT    NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS categories_parent_idx ON categories (parent_id, sort_order);

ALTER TABLE products
    ADD COLUMN IF NOT EXISTS category_id BIGINT REFERENCES categories (id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS products_category_idx ON products (category_id);

-- ─────────────────────────────────────────────────────────────────────────
-- Seed: 5 parent groups + 17 leaves. Slug is the public identifier; the
-- bigserial id is internal, so children are linked via parent slug lookup.
-- ─────────────────────────────────────────────────────────────────────────

INSERT INTO categories (slug, name, parent_id, sort_order)
VALUES ('priroda', 'Природа', NULL, 10),
       ('prazdniki-grp', 'Праздники', NULL, 20),
       ('tematika', 'Тематика', NULL, 30),
       ('gorod-byt', 'Город и быт', NULL, 40),
       ('eda-napitki', 'Еда и напитки', NULL, 50)
ON CONFLICT (slug) DO NOTHING;

INSERT INTO categories (slug, name, parent_id, sort_order)
VALUES
    -- Природа
    ('pticy', 'Птицы', (SELECT id FROM categories WHERE slug = 'priroda'), 10),
    ('zhivotnyy-mir', 'Животный мир', (SELECT id FROM categories WHERE slug = 'priroda'), 20),
    ('cvety', 'Цветы', (SELECT id FROM categories WHERE slug = 'priroda'), 30),
    ('venochki', 'Веночки разные', (SELECT id FROM categories WHERE slug = 'priroda'), 40),
    ('uzory-ornamenty', 'Узоры. Орнаменты. Фоны.', (SELECT id FROM categories WHERE slug = 'priroda'), 50),
    -- Праздники
    ('novyy-god', 'Новый год и Рождество', (SELECT id FROM categories WHERE slug = 'prazdniki-grp'), 10),
    ('prazdniki', 'Праздники', (SELECT id FROM categories WHERE slug = 'prazdniki-grp'), 20),
    ('gnomy', 'Гномы', (SELECT id FROM categories WHERE slug = 'prazdniki-grp'), 30),
    -- Тематика
    ('tekst-bukvy', 'Текст. Буквы', (SELECT id FROM categories WHERE slug = 'tematika'), 10),
    ('muzyka-goroskop', 'Музыка. Гороскоп.', (SELECT id FROM categories WHERE slug = 'tematika'), 20),
    ('lyudi', 'Люди', (SELECT id FROM categories WHERE slug = 'tematika'), 30),
    ('transport-igry', 'Транспорт. Карты. Игры. Мужская тема.', (SELECT id FROM categories WHERE slug = 'tematika'), 40),
    -- Город и быт
    ('posuda-odezhda', 'Посуда. Одежда. Обувь. Предметы быта.', (SELECT id FROM categories WHERE slug = 'gorod-byt'), 10),
    ('goroda-peyzazh', 'Города. Пейзаж. Дома.', (SELECT id FROM categories WHERE slug = 'gorod-byt'), 20),
    -- Еда и напитки
    ('produktovaya-tema', 'Продуктовая тема', (SELECT id FROM categories WHERE slug = 'eda-napitki'), 10),
    ('frukty-yagody', 'Фрукты и ягоды', (SELECT id FROM categories WHERE slug = 'eda-napitki'), 20),
    ('vino-vinograd', 'Вино и виноград', (SELECT id FROM categories WHERE slug = 'eda-napitki'), 30)
ON CONFLICT (slug) DO NOTHING;

-- ─────────────────────────────────────────────────────────────────────────
-- Demo products. Slug is the public identifier used by the FE.
-- Атрибуты товаров живут в отдельной таблице (см. 0003).
-- ─────────────────────────────────────────────────────────────────────────

INSERT INTO products (slug, name, category_id)
VALUES ('roses-napkin', 'Салфетка с розами',
        (SELECT id FROM categories WHERE slug = 'cvety')),
       ('daisies-napkin', 'Салфетка с ромашками',
        (SELECT id FROM categories WHERE slug = 'cvety')),
       ('butterflies-napkin', 'Салфетка с бабочками',
        (SELECT id FROM categories WHERE slug = 'zhivotnyy-mir')),
       ('leaves-napkin', 'Салфетка с листьями',
        (SELECT id FROM categories WHERE slug = 'uzory-ornamenty')),
       ('gnomes-napkin', 'Гномы новогодние',
        (SELECT id FROM categories WHERE slug = 'gnomy'))
ON CONFLICT (slug) DO NOTHING;