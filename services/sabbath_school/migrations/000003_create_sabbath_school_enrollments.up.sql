CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS sabbath_school_enrollments (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sabbath_school_id UUID NOT NULL,
    member_id         UUID NOT NULL,
    enrolled_at       DATE NOT NULL DEFAULT CURRENT_DATE,
    unenrolled_at     DATE CHECK (unenrolled_at IS NULL OR unenrolled_at >= enrolled_at),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by        UUID NOT NULL,
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by        UUID,
    deleted_at        TIMESTAMPTZ,
    deleted_by        UUID
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_sabbath_school_enrollments_active
    ON sabbath_school_enrollments (sabbath_school_id, member_id)
    WHERE unenrolled_at IS NULL AND deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_sabbath_school_enrollments_school
    ON sabbath_school_enrollments (sabbath_school_id);

CREATE INDEX IF NOT EXISTS idx_sabbath_school_enrollments_member
    ON sabbath_school_enrollments (member_id);
