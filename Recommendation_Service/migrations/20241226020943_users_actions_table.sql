-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_actions (
                              id SERIAL PRIMARY KEY,
                              user_id INT NOT NULL,
                              product_id INT NOT NULL,
                              name VARCHAR(255) NOT NULL,
                              category VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_actions;
-- +goose StatementEnd
