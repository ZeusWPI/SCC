-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    academic_year TEXT NOT NULL,
    location TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event;
-- +goose StatementEnd
