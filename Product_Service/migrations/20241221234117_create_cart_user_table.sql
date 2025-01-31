-- +goose Up
-- +goose StatementBegin
CREATE TABLE cart (
                      id SERIAL PRIMARY KEY,          -- Unique identifier for the cart item
                      user_id VARCHAR(255) NOT NULL,  -- User ID
                      product_id INT NOT NULL,        -- Product ID
                      quantity INT NOT NULL,
                      FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cart;
-- +goose StatementEnd
