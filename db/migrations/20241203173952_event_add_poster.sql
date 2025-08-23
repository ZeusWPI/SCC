-- +goose Up
-- +goose StatementBegin
ALTER TABLE event
ADD COLUMN poster BYTES;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE event
DROP COLUMN poster;
-- +goose StatementEnd
