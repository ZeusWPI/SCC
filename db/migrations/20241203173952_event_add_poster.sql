-- +goose Up
-- +goose StatementBegin
ALTER TABLE event
add COLUMN poster BLOB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE event
DROP column poster;
-- +goose StatementEnd
