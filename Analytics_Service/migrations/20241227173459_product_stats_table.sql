-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product_stats (
                                             id SERIAL PRIMARY KEY,
                                             product_id INT NOT NULL,
                                             cart_add_count INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE product_stats;
-- +goose StatementEnd
