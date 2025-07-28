-- +goose Up
ALTER TABLE orders
ADD COLUMN status TEXT NOT NULL DEFAULT 'pending';

-- +goose Down
ALTER TABLE orders
DROP COLUMN status;
