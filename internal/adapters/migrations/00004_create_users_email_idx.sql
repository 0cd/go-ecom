-- +goose no transaction

-- +goose Up
CREATE INDEX CONCURRENTLY ON users(email);

-- +goose Down
DROP INDEX users_email_idx;
