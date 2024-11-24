-- +goose Up
-- +goose StatementBegin
ALTER TABLE season
ADD COLUMN current BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE season
DROP COLUMN current;
-- +goose StatementEnd
