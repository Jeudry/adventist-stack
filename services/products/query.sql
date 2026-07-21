-- name: CreateProduct :one
INSERT INTO products (name, sku, description, brand, release_date, status)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: ListProducts :many
SELECT * FROM products
WHERE (
    @search::text = ''
    OR name ILIKE ('%' || @search || '%')
    OR sku  ILIKE ('%' || @search || '%')
)
ORDER BY created_at DESC
LIMIT @row_limit OFFSET @row_offset;

-- name: CountProducts :one
SELECT count(*) FROM products
WHERE (
    @search::text = ''
    OR name ILIKE '%' || @search || '%'
    OR sku  ILIKE '%' || @search || '%'
);

-- name: UpdateProduct :one
UPDATE products SET
    name = $2, sku = $3, description = $4, brand = $5,
    release_date = $6, status = $7, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :execrows
DELETE FROM products WHERE id = $1;
