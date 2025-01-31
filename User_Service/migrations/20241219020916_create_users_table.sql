-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS goose_db_version (
                                                id SERIAL PRIMARY KEY,
                                                version_id BIGINT NOT NULL,
                                                is_applied BOOLEAN NOT NULL,
                                                tstamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       email VARCHAR(255) NOT NULL UNIQUE,
                       password_hash VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
