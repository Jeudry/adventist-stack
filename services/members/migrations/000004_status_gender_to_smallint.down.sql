ALTER TABLE members
    ALTER COLUMN gender TYPE VARCHAR(20)
    USING (CASE gender
        WHEN 1 THEN 'M'
        WHEN 2 THEN 'F'
        ELSE 'M'
    END);
ALTER TABLE members ADD CONSTRAINT members_gender_check CHECK (char_length(gender) > 0);

ALTER TABLE members ALTER COLUMN status DROP DEFAULT;
ALTER TABLE members
    ALTER COLUMN status TYPE VARCHAR(20)
    USING (CASE status
        WHEN 1 THEN 'active'
        WHEN 2 THEN 'inactive'
        WHEN 3 THEN 'visitor'
        ELSE 'active'
    END);
ALTER TABLE members ALTER COLUMN status SET DEFAULT 'active';
ALTER TABLE members ADD CONSTRAINT members_status_check
    CHECK (status IN ('active', 'inactive', 'visitor'));
