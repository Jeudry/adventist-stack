-- name: CreateMember :one
INSERT INTO members (first_name, last_name, email, phone, gender, address, birth_date, baptism_date, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetMemberByID :one
SELECT * FROM members WHERE id = $1;

-- name: ListMembers :many
SELECT * FROM members WHERE (
    @search::TEXT = ''
    OR first_name ILIKE ('%' || @search || '%')
    OR last_name ILIKE ('%' || @search || '%')
    OR email ILIKE ('%' || @search || '%')
) ORDER BY created_at DESC
LIMIT @row_limit OFFSET @row_offset;

-- name: CountMembers :one
SELECT count(*) FROM members
WHERE (
    @search::text = ''
    OR first_name ILIKE ('%' || @search || '%')
    OR last_name ILIKE ('%' || @search || '%')
    OR email ILIKE ('%' || @search || '%')
);

-- name: UpdateMember :one
UPDATE members SET 
    first_name = $2,
    last_name = $3,
    email = $4,
    phone = $5,
    gender = $6,
    address = $7,
    birth_date = $8,
    baptism_date = $9,
    status = $10,
    updated_at = NOW()
WHERE id = $1 
RETURNING *;

-- name: DeleteMember :execrows
DELETE FROM members WHERE id = $1;    