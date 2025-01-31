-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS goose_db_version (
                                                id SERIAL PRIMARY KEY,
                                                version_id BIGINT NOT NULL,
                                                is_applied BOOLEAN NOT NULL,
                                                tstamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS user_actions (
                                            id SERIAL PRIMARY KEY,
                                            user_id INT NOT NULL,
                                            product_id INT NOT NULL,
                                            action_type VARCHAR(255) NOT NULL,
                                            product_name VARCHAR(255) NOT NULL,
                                            category VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_actions;
-- +goose StatementEnd
