ALTER TABLE prayers ALTER COLUMN status DROP DEFAULT;
ALTER TABLE prayers
    ALTER COLUMN status TYPE VARCHAR(20)
    USING (CASE status
        WHEN 1 THEN 'pending'
        WHEN 2 THEN 'answered'
        WHEN 3 THEN 'archived'
        ELSE 'pending'
    END);
ALTER TABLE prayers ALTER COLUMN status SET DEFAULT 'pending';
ALTER TABLE prayers ADD CONSTRAINT prayers_status_check
    CHECK (status IN ('pending', 'answered', 'archived'));
