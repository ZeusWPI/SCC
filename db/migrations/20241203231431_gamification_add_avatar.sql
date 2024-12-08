-- +goose Up
-- +goose StatementBegin
ALTER TABLE gamification
DROP COLUMN avatar;

ALTER TABLE gamification
ADD COLUMN avatar BYTEA;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE gamification
DROP COLUMN avatar;

ALTER TABLE gamification
ADD COLUMN avatar VARCHAR(255) NOT NULL;
-- +goose StatementEnd
