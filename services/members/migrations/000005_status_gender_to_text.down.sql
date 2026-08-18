ALTER TABLE members ALTER COLUMN status DROP DEFAULT;
ALTER TABLE members DROP CONSTRAINT IF EXISTS members_status_check;
ALTER TABLE members
    ALTER COLUMN status TYPE SMALLINT
    USING (CASE status
        WHEN 'active' THEN 1
        WHEN 'inactive' THEN 2
        WHEN 'visitor' THEN 3
        ELSE 1
    END);
ALTER TABLE members ALTER COLUMN status SET DEFAULT 1;

ALTER TABLE members DROP CONSTRAINT IF EXISTS members_gender_check;
ALTER TABLE members
    ALTER COLUMN gender TYPE SMALLINT
    USING (CASE upper(gender)
        WHEN 'M' THEN 1
        WHEN 'F' THEN 2
        ELSE 1
    END);
