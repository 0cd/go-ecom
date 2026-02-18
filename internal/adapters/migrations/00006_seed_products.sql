-- +goose Up
-- +goose StatementBegin
INSERT INTO products (name, price_in_cents, quantity)
VALUES
  ('Wireless Mouse', 2599, 50),
  ('Bluetooth Headphones', 4999, 30),
  ('USB-C Charger', 1999, 100),
  ('Mechanical Keyboard', 8999, 20),
  ('Webcam 1080p', 3499, 40),
  ('Laptop Stand', 2999, 25),
  ('Portable SSD 1TB', 12999, 15),
  ('Smartwatch', 15999, 10),
  ('Gaming Chair', 22999, 8),
  ('Desk Lamp', 1799, 60);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM products
WHERE name IN (
  'Wireless Mouse',
  'Bluetooth Headphones',
  'USB-C Charger',
  'Mechanical Keyboard',
  'Webcam 1080p',
  'Laptop Stand',
  'Portable SSD 1TB',
  'Smartwatch',
  'Gaming Chair',
  'Desk Lamp'
);
-- +goose StatementEnd
