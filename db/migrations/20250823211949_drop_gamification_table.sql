-- +goose Up
-- +goose StatementBegin
DROP TABLE gamification;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE gamification (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    score INTEGER NOT NULL,
    avatar BYTEA
);

-- +goose StatementEnd
