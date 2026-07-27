-- name: CreatePayer :one
INSERT INTO prayers (
    id,
    title,
    description,
    author_name,
    is_anonymous,
    status,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetPrayerByID :one
SELECT * FROM prayers WHERE id = $1;

-- name: ListPrayers :many
SELECt * FROM prayers WHERE (
    @search::TEXT = ''
    OR title ILIKE ('%' || @search || '%')
    OR description ILIKE ('%' || @search || '%')
    OR author_name ILIKE ('%' || @search || '%')
) ORDER BY created_at DESC 
LIMIT @row_limit OFFSET @row_offset;

-- name: Count :one
SELECt count(*) FROM prayers WHERE (
    @search::TEXT = ''
    OR title ILIKE ('%' || @search || '%')
    OR description ILIKE ('%' || @search || '%')
    OR author_name ILIKE ('%' || @search || '%')
);

-- name: UpdatePrayer :one
UPDATE prayers SET 
    title = $2,
    description = $3,
    author_name = $4,
    is_anonymous = $5,
    status = $6,
    updated_at = NOW(),
    updated_by = $7
WHERE id = $1
RETURNING *;

-- name: DeletePrayer :exec
DELETE FROM prayers WHERE id = $1;