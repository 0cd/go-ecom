-- name: ListProducts :many
SELECT * FROM products ORDER BY id;

-- name: FindProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: CreateProduct :one
INSERT INTO products (
  name, price_in_cents, quantity
) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET
  name = coalesce(sqlc.narg('name'), name),
  price_in_cents = coalesce(sqlc.narg('price_in_cents'), price_in_cents),
  quantity = coalesce(sqlc.narg('quantity'), quantity)
WHERE id = $1 RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: CreateOrder :one
INSERT INTO orders (
  customer_id
) VALUES ($1) RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (
  order_id, product_id, quantity, price_in_cents
) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: FindOrderByID :many
SELECT orders.id, orders.customer_id, orders.created_at, order_items.product_id, order_items.quantity, order_items.price_in_cents FROM orders
INNER JOIN order_items ON order_items.order_id = orders.id
WHERE orders.id = $1
ORDER BY order_items.product_id;

-- name: CreateUser :one
INSERT INTO users (
  email, password_hash
) VALUES ($1, $2) RETURNING *;

-- name: ListUsers :many
SELECT id, email, verified, is_admin, updated_at, created_at
FROM users
ORDER BY id;

-- name: SearchUsers :many
SELECT id, email, verified, is_admin, updated_at, created_at
FROM users
WHERE email ILIKE '%' || $1 || '%'
ORDER BY id;

-- name: FindUserByID :one
SELECT id, email, verified, is_admin, updated_at, created_at
FROM users
WHERE id = $1;

-- name: FindUserByEmail :one
SELECT id, email, verified, is_admin, updated_at, created_at
FROM users
WHERE email = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2, updated_at = now()
WHERE id = $1;

-- name: UpdateUserEmail :exec
UPDATE users
SET email = $2, verified = false, updated_at = now()
WHERE id = $1;

-- name: VerifyUser :exec
UPDATE users
SET verified = true, updated_at = now()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;