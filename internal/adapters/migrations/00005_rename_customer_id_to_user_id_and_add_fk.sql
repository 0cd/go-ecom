-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
RENAME COLUMN customer_id TO user_id;

ALTER TABLE orders
ADD CONSTRAINT fk_orders_user_id
FOREIGN KEY (user_id) REFERENCES users(id)
ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
DROP CONSTRAINT IF EXISTS fk_orders_user_id;

ALTER TABLE orders
RENAME COLUMN user_id TO customer_id;
-- +goose StatementEnd
