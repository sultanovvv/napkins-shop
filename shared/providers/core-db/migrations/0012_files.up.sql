BEGIN;

-- Оригинал и каждый производный размер имеют собственную запись и UUID.
CREATE TABLE files
(
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status             TEXT NOT NULL DEFAULT 'pending',
    s3_key             TEXT UNIQUE,
    mime_type          TEXT,
    width              INTEGER,
    height             INTEGER,
    byte_size          BIGINT,
    unreferenced_since TIMESTAMPTZ DEFAULT now(),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT files_status_valid CHECK (status IN ('pending', 'processing', 'ready', 'failed', 'deleting')),
    CONSTRAINT files_key_not_blank CHECK (s3_key IS NULL OR length(btrim(s3_key)) > 0),
    CONSTRAINT files_mime_not_blank CHECK (mime_type IS NULL OR length(btrim(mime_type)) > 0),
    CONSTRAINT files_width_positive CHECK (width > 0),
    CONSTRAINT files_height_positive CHECK (height > 0),
    CONSTRAINT files_size_positive CHECK (byte_size > 0),
    CONSTRAINT files_ready_complete CHECK (
        status <> 'ready' OR (
            s3_key IS NOT NULL AND mime_type IS NOT NULL
            AND width IS NOT NULL AND height IS NOT NULL AND byte_size IS NOT NULL
        )
    )
);

COMMENT ON COLUMN files.unreferenced_since IS
    'Начало срока хранения без товара; NULL при наличии связей. Обновляется файловым сервисом';

CREATE INDEX files_cleanup_idx ON files (unreferenced_since)
    WHERE unreferenced_since IS NOT NULL;

CREATE TABLE file_variants
(
    source_file_id UUID NOT NULL REFERENCES files (id) ON DELETE RESTRICT,
    preset        TEXT NOT NULL,
    file_id       UUID NOT NULL UNIQUE REFERENCES files (id) ON DELETE RESTRICT,
    profile_version INTEGER NOT NULL DEFAULT 1,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (source_file_id, preset),
    CONSTRAINT file_variants_different_files CHECK (source_file_id <> file_id),
    CONSTRAINT file_variants_preset_not_blank CHECK (length(btrim(preset)) > 0),
    CONSTRAINT file_variants_profile_positive CHECK (profile_version > 0)
);

-- Расширяем старую таблицу без удаления ключей, порядка и признака главной.
-- file_id заполняется после чтения оригиналов и подготовки файловым сервисом.
ALTER TABLE product_images
    ADD COLUMN file_id UUID REFERENCES files (id) ON DELETE RESTRICT,
    ADD CONSTRAINT product_images_product_file_unique UNIQUE (product_id, file_id);

CREATE INDEX product_images_file_idx ON product_images (file_id) WHERE file_id IS NOT NULL;

COMMIT;
