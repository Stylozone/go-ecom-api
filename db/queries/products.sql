-- name: CreateProduct :one
INSERT INTO products (name, description, price, quantity, image_url)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY id DESC
LIMIT $1 OFFSET $2;

-- name: UpdateProductQuantity :exec
UPDATE products
SET quantity = $2
WHERE id = $1;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;
