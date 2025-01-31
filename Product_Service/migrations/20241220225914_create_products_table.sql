-- +goose Up
-- +goose StatementBegin
CREATE TABLE products (
                          id SERIAL PRIMARY KEY,          -- Unique product identifier
                          name VARCHAR(255) NOT NULL,     -- Product name
                          category VARCHAR(255) NOT NULL, -- Product category
                          price INT NOT NULL              -- Product price
);
CREATE INDEX idx_products_category ON products(category);
INSERT INTO products (name, category, price) VALUES
                                                 ('White T-shirt', 'Clothing', 20),
                                                 ('Black Jeans', 'Clothing', 50),
                                                 ('Winter Jacket', 'Clothing', 150),
                                                 ('Summer Shorts', 'Clothing', 30),
                                                 ('Warm Sweater', 'Clothing', 80);
INSERT INTO products (name, category, price) VALUES
                                                 ('ASUS Laptop', 'Electronics', 1200),
                                                 ('iPhone 14', 'Electronics', 900),
                                                 ('Samsung Tablet', 'Electronics', 600),
                                                 ('Sony Headphones', 'Electronics', 200),
                                                 ('LG TV', 'Electronics', 1500);
INSERT INTO products (name, category, price) VALUES
                                                 ('Ray-Ban Sunglasses', 'Glasses', 150),
                                                 ('Reading Glasses', 'Glasses', 50),
                                                 ('Swimming Goggles', 'Glasses', 30),
                                                 ('Driving Glasses', 'Glasses', 70),
                                                 ('UV Protection Glasses', 'Glasses', 100);
-- +goose StatementEnd

-- +goose Down
    DROP INDEX IF EXISTS idx_products_category;
    DROP TABLE products;
-- +goose StatementBegin

-- +goose StatementEnd
