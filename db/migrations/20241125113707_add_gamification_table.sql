-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS gamification (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    score INTEGER NOT NULL,
    avatar VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS gamification;
-- +goose StatementEnd
