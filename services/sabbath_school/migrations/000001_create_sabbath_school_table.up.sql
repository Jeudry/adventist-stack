CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS sabbath_school (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name           VARCHAR(128) NOT NULL CHECK(char_length(name) >= 3),
    teacher_id     UUID,
    location       VARCHAR(128) CHECK(location IS NULL OR char_length(location) >= 2),
    target_min_age INTEGER CHECK(target_min_age IS NULL OR (target_min_age >= 0 AND target_min_age <= 120)),
    target_max_age INTEGER CHECK(target_max_age IS NULL OR (target_max_age >= target_min_age AND target_max_age <= 120)),
    status         INTEGER NOT NULL DEFAULT 1 CHECK(status >= 0 AND status <= 10),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by     UUID NOT NULL,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by     UUID,
    deleted_at     TIMESTAMPTZ,
    deleted_by     UUID
);

CREATE INDEX IF NOT EXISTS idx_sabbath_school_status ON sabbath_school (status);
CREATE INDEX IF NOT EXISTS idx_sabbath_school_teacher_id ON sabbath_school (teacher_id);
CREATE INDEX IF NOT EXISTS idx_sabbath_school_created_at ON sabbath_school (created_at DESC);