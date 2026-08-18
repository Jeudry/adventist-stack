ALTER TABLE prayers ALTER COLUMN status DROP DEFAULT;
ALTER TABLE prayers DROP CONSTRAINT IF EXISTS prayers_status_check;
ALTER TABLE prayers
    ALTER COLUMN status TYPE INTEGER
    USING (CASE status
        WHEN 'pending' THEN 1
        WHEN 'answered' THEN 2
        WHEN 'archived' THEN 3
        ELSE 0
    END);
ALTER TABLE prayers ALTER COLUMN status SET DEFAULT 0;
