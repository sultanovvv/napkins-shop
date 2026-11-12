BEGIN;

-- Старый CHECK проверял заполненность, но не соответствие типу справочника.
-- Несогласованные данные требуют исправления, а не автоматического удаления.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM product_attribute_values AS v
        JOIN attributes AS a ON a.id = v.attribute_id
        WHERE NOT (
            (a.value_type = 'string' AND v.value_text IS NOT NULL
                AND length(btrim(v.value_text)) > 0 AND v.value_int IS NULL)
            OR (a.value_type = 'int' AND v.value_int IS NOT NULL AND v.value_text IS NULL)
        )
    ) THEN
        RAISE EXCEPTION 'Existing product attribute values do not match their definitions'
            USING ERRCODE = '23514';
    END IF;
END;
$$;

ALTER TABLE attributes
    ADD COLUMN unit TEXT,
    ADD CONSTRAINT attributes_name_not_blank CHECK (length(btrim(name)) > 0),
    ADD CONSTRAINT attributes_slug_not_blank CHECK (length(btrim(slug)) > 0),
    ADD CONSTRAINT attributes_unit_valid CHECK (
        unit IS NULL OR (length(btrim(unit)) > 0 AND value_type IN ('int', 'decimal'))
    ),
    ADD CONSTRAINT attributes_id_type_unique UNIQUE (id, value_type);

ALTER TABLE product_attribute_values
    ADD COLUMN value_type attribute_value_type,
    ADD COLUMN value_decimal NUMERIC(18,6),
    ADD COLUMN value_bool BOOLEAN;

UPDATE product_attribute_values AS v
SET value_type = a.value_type
FROM attributes AS a
WHERE a.id = v.attribute_id;

ALTER TABLE product_attribute_values
    ALTER COLUMN value_type SET NOT NULL,
    DROP CONSTRAINT product_attribute_values_check,
    DROP CONSTRAINT product_attribute_values_attribute_id_fkey,
    ADD CONSTRAINT product_attribute_values_definition_fk
        FOREIGN KEY (attribute_id, value_type) REFERENCES attributes (id, value_type)
        ON DELETE RESTRICT ON UPDATE RESTRICT,
    ADD CONSTRAINT product_attribute_values_typed_value CHECK (
        CASE value_type
            WHEN 'string' THEN value_text IS NOT NULL
                AND value_int IS NULL AND value_decimal IS NULL AND value_bool IS NULL
            WHEN 'int' THEN value_int IS NOT NULL
                AND value_text IS NULL AND value_decimal IS NULL AND value_bool IS NULL
            WHEN 'decimal' THEN value_decimal IS NOT NULL
                AND value_text IS NULL AND value_int IS NULL AND value_bool IS NULL
            WHEN 'boolean' THEN value_bool IS NOT NULL
                AND value_text IS NULL AND value_int IS NULL AND value_decimal IS NULL
            ELSE FALSE
        END
    ),
    ADD CONSTRAINT product_attribute_values_text_not_blank
        CHECK (value_text IS NULL OR length(btrim(value_text)) > 0),
    ADD CONSTRAINT product_attribute_values_decimal_finite CHECK (
        value_decimal IS NULL
        OR value_decimal NOT IN ('NaN'::numeric, 'Infinity'::numeric, '-Infinity'::numeric)
    );

COMMIT;
