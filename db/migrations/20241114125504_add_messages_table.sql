-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS message (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    ip TEXT NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS message;
-- +goose StatementEnd
