ALTER TABLE products ALTER COLUMN status DROP DEFAULT;
ALTER TABLE products DROP CONSTRAINT IF EXISTS products_status_check;

ALTER TABLE products
    ALTER COLUMN status TYPE SMALLINT
    USING (CASE status
        WHEN 'active' THEN 1
        WHEN 'inactive' THEN 2
        WHEN 'discontinued' THEN 3
        ELSE 1
    END);

ALTER TABLE products ALTER COLUMN status SET DEFAULT 1;
