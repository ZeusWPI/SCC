-- +goose Up
-- +goose StatementBegin
ALTER TABLE event
ADD COLUMN poster BLOB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE event
DROP COLUMN poster;
-- +goose StatementEnd
