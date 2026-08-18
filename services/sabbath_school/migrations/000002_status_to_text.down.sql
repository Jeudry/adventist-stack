ALTER TABLE sabbath_school ALTER COLUMN status DROP DEFAULT;
ALTER TABLE sabbath_school DROP CONSTRAINT IF EXISTS sabbath_school_status_check;
ALTER TABLE sabbath_school
    ALTER COLUMN status TYPE INTEGER
    USING (CASE status
        WHEN 'active' THEN 1
        WHEN 'inactive' THEN 2
        ELSE 1
    END);
ALTER TABLE sabbath_school ALTER COLUMN status SET DEFAULT 1;
ALTER TABLE sabbath_school ADD CONSTRAINT sabbath_school_status_check
    CHECK (status >= 0 AND status <= 10);
