CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS prayers (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title        VARCHAR(128) NOT NULL CHECK (char_length(title) >= 3),
    description  VARCHAR(2056) NOT NULL CHECK (char_length(description) >= 5),
    author_name VARCHAR(256) CHECK (author_name IS NULL OR char_length(author_name) >= 3),
    is_anonymous BOOLEAN NOT NULL DEFAULT FALSE,
    status       INTEGER NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    created_by   UUID         NOT NULL,
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_by   UUID,
    deleted_at   TIMESTAMPTZ,
    deleted_by   UUID
);

CREATE INDEX IF NOT EXISTS idx_prayers_status ON prayers (status);
CREATE INDEX IF NOT EXISTS idx_prayers_created_at ON prayers (created_at DESC);