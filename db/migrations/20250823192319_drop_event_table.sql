-- +goose Up
-- +goose StatementBegin
DROP TABLE event;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE event (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    academic_year TEXT NOT NULL,
    location TEXT NOT NULL,
    poster BYTES
);
-- +goose StatementEnd
