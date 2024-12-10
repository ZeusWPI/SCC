-- +goose Up
-- +goose StatementBegin
ALTER TABLE scan
ADD COLUMN scan_id INTEGER NOT NULL DEFAULT 0;

ALTER TABLE scan
ALTER COLUMN scan_id
DROP DEFAULT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE scan
DROP COLUMN scan_id;
-- +goose StatementEnd
