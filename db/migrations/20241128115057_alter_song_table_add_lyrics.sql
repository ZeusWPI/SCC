-- +goose Up
-- +goose StatementBegin
ALTER TABLE song
ADD COLUMN album TEXT NOT NULL;

ALTER TABLE song
ADD COLUMN lyrics_type TEXT;

ALTER TABLE song
ADD COLUMN lyrics TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE song
DROP COLUMN isrc_id;

ALTER TABLE song
DROP COLUMN lyrics;

ALTER TABLE song
DROP COLUMN common_id;
-- +goose StatementEnd
