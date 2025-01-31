-- +goose Up
-- +goose StatementBegin
CREATE TABLE products (
                          id SERIAL PRIMARY KEY,
                          product_id INT NOT NULL,
                          name VARCHAR(255) NOT NULL,
                          category VARCHAR(255) NOT NULL
);
CREATE INDEX idx_products_category ON products(category);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE products;
-- +goose StatementEnd
