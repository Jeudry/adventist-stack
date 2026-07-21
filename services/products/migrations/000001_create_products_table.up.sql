CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS products (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name         VARCHAR(150) NOT NULL CHECK (char_length(name) > 0),
    sku          VARCHAR(50)  NOT NULL UNIQUE,
    description  VARCHAR(500),
    brand        VARCHAR(100),
    release_date DATE,
    status       VARCHAR(20)  NOT NULL DEFAULT 'active'
                 CHECK (status IN ('active','inactive','discontinued')),
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_products_sku ON products (sku);
