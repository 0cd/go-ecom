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
  user_id
) VALUES ($1) RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (
  order_id, product_id, quantity, price_in_cents
) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: FindOrderByID :many
SELECT o.id, o.user_id, o.created_at, oi.product_id, oi.quantity, oi.price_in_cents FROM orders o
INNER JOIN order_items oi ON oi.order_id = o.id
WHERE o.id = $1
ORDER BY oi.product_id;

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
SELECT id, email, verified, password_hash, is_admin, updated_at, created_at
FROM users
WHERE id = $1;

-- name: FindUserByEmail :one
SELECT id, email, verified, password_hash, is_admin, updated_at, created_at
FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET
  email = coalesce(sqlc.narg('email'), email),
  password_hash = coalesce(sqlc.narg('password_hash'), password_hash),
  verified = coalesce(sqlc.narg('verified'), verified),
  updated_at = now()
WHERE id = $1 RETURNING *;


-- name: VerifyUser :exec
UPDATE users
SET verified = true, updated_at = now()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;