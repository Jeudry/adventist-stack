ALTER TABLE users ALTER COLUMN role DROP DEFAULT;
ALTER TABLE users
    ALTER COLUMN role TYPE VARCHAR(20)
    USING (CASE role
        WHEN 1 THEN 'admin'
        WHEN 2 THEN 'member'
        ELSE 'member'
    END);
ALTER TABLE users ALTER COLUMN role SET DEFAULT 'member';
ALTER TABLE users ADD CONSTRAINT users_role_check
    CHECK (role IN ('admin', 'member'));
