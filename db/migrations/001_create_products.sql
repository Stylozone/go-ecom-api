-- +goose Up
CREATE TABLE IF NOT EXISTS products (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  price INT NOT NULL,
  quantity INT NOT NULL,
  created_at TIMESTAMP DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS products;
