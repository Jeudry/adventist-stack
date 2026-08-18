ALTER TABLE sabbath_school ALTER COLUMN status DROP DEFAULT;
ALTER TABLE sabbath_school DROP CONSTRAINT IF EXISTS sabbath_school_status_check;
ALTER TABLE sabbath_school
    ALTER COLUMN status TYPE VARCHAR(20)
    USING (CASE status
        WHEN 1 THEN 'active'
        WHEN 2 THEN 'inactive'
        ELSE 'active'
    END);
ALTER TABLE sabbath_school ALTER COLUMN status SET DEFAULT 'active';
ALTER TABLE sabbath_school ADD CONSTRAINT sabbath_school_status_check
    CHECK (status IN ('active', 'inactive'));
