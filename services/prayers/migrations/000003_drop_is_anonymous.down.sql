ALTER TABLE prayers ADD COLUMN is_anonymous BOOLEAN NOT NULL DEFAULT FALSE;
UPDATE prayers SET is_anonymous = (author_name IS NULL);
