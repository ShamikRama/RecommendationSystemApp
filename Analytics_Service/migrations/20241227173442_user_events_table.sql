-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_events (
                                           id SERIAL PRIMARY KEY,
                                           user_id INT NOT NULL,
                                           event_type VARCHAR(255) NOT NULL,
                                           email VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_events;
-- +goose StatementEnd
