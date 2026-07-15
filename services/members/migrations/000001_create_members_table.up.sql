CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS members (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name   VARCHAR(100) NOT NULL CHECK (char_length(first_name) > 0),
    last_name    VARCHAR(100) NOT NULL CHECK (char_length(last_name) > 0),
    email        VARCHAR(254) CHECK (email ~ '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    phone        VARCHAR(30)  CHECK (phone ~ '^[+0-9 ()-]{7,30}$'),
    gender       VARCHAR(20)  NOT NULL CHECK (char_length(gender) > 0),
    status       VARCHAR(20)  NOT NULL DEFAULT 'active'
                 CHECK (status IN ('active','inactive','visitor')),
    birth_date   DATE,
    baptism_date DATE,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_members_email ON members (email);
CREATE INDEX IF NOT EXISTS idx_members_phone ON members (phone);
