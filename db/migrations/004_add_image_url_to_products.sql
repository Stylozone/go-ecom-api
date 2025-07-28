-- +goose Up
ALTER TABLE products ADD COLUMN image_url TEXT;

-- +goose Down
ALTER TABLE products DROP COLUMN image_url;
