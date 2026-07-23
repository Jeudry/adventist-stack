ALTER TABLE products ALTER COLUMN status DROP DEFAULT;

ALTER TABLE products
    ALTER COLUMN status TYPE VARCHAR(20)
    USING (CASE status
        WHEN 1 THEN 'active'
        WHEN 2 THEN 'inactive'
        WHEN 3 THEN 'discontinued'
        ELSE 'active'
    END);

ALTER TABLE products ALTER COLUMN status SET DEFAULT 'active';
ALTER TABLE products ADD CONSTRAINT products_status_check
    CHECK (status IN ('active', 'inactive', 'discontinued'));
