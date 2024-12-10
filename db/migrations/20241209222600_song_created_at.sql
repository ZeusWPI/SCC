-- +goose Up
-- +goose StatementBegin
ALTER TABLE song_history
ALTER COLUMN created_at
SET DEFAULT (NOW() - INTERVAL '1 second');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE song_history
ALTER COLUMN created_at
SET DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd
