-- +goose Up
-- +goose StatementBegin
ALTER TABLE spotify
DROP COLUMN created_at;

ALTER TABLE spotify RENAME TO song;

CREATE TABLE IF NOT EXISTS song_history (
    id INTEGER PRIMARY KEY,
    song_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(song_id) REFERENCES song(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS song_history;

ALTER TABLE song RENAME TO spotify;

ALTER TABLE spotify
ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd
