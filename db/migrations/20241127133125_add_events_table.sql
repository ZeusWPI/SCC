-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    date TIMESTAMP NOT NULL,
    academic_year TEXT NOT NULL,
    location TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event;
-- +goose StatementEnd
