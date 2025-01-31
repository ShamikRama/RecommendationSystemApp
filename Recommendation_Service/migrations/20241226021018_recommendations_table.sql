-- +goose Up
-- +goose StatementBegin
CREATE TABLE recommendations (
                                 id SERIAL PRIMARY KEY,
                                 user_id INT NOT NULL,
                                 product_id INT NOT NULL,
                                 name VARCHAR(255) NOT NULL,
                                 category VARCHAR(255) NOT NULL,
                                 timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE recommendations;
-- +goose StatementEnd
