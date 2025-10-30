-- +goose Up
-- +goose StatementBegin
ALTER TABLE song_artist
DROP COLUMN followers;

ALTER TABLE song_artist
DROP COLUMN popularity;

CREATE TYPE lyrics_type_enum AS ENUM ('plain', 'synced', 'instrumental', 'missing');

UPDATE song
SET lyrics_type = 'missing'
WHERE lyrics_type = '';

ALTER TABLE song
ADD COLUMN lyrics_type_new lyrics_type_enum NOT NULL DEFAULT 'missing';

UPDATE song
SET lyrics_type_new = lyrics_type::lyrics_type_enum;

ALTER TABLE song
DROP COLUMN lyrics_type;

ALTER TABLE song
RENAME COLUMN lyrics_type_new TO lyrics_type;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE song
ADD COLUMN lyrics_type_text TEXT;

UPDATE song
SET lyrics_type_text = lyrics_type::text;

ALTER TABLE song
DROP COLUMN lyrics_type;

ALTER TABLE song
RENAME COLUMN lyrics_type_text TO lyrics_type;

DROP TYPE lyrics_type_enum;

ALTER TABLE song_artist
ADD COLUMN popularity INT;

ALTER TABLE song_artist
ADD COLUMN followers INT;
-- +goose StatementEnd
