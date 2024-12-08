-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS spotify (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    artists TEXT NOT NULL,
    spotify_id TEXT NOT NULL,
    duration_ms INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS spotify;
-- +goose StatementEnd
