-- name: CreateOrder :one
INSERT INTO orders (user_id, total_price)
VALUES ($1, $2)
RETURNING *;

-- name: CreateOrderItem :exec
INSERT INTO order_items (order_id, product_id, quantity, price)
VALUES ($1, $2, $3, $4);

-- name: GetOrdersByUser :many
SELECT o.id AS order_id, o.total_price, o.created_at,
       oi.product_id, oi.quantity, oi.price
FROM orders o
JOIN order_items oi ON o.id = oi.order_id
WHERE o.user_id = $1
ORDER BY o.created_at DESC;

-- name: UpdateOrderStatus :exec
UPDATE orders
SET status = $2
WHERE id = $1;
