-- name: CreateUser :one
INSERT INTO users (email, name, password_hash, role)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: ExistsUserByEmail :one
SELECT EXISTS (SELECT 1 FROM users WHERE email = $1);
