-- +goose Up
-- +goose StatementBegin
CREATE TYPE TAP_CATEGORY AS ENUM ('soft', 'mate', 'beer', 'food', 'unknown');

ALTER TABLE tap
ADD COLUMN new_category TAP_CATEGORY NOT NULL DEFAULT 'unknown';

UPDATE tap
SET new_category = CASE 
    WHEN LOWER(category) IN ('soft', 'mate', 'beer', 'food', 'unknown')
         THEN LOWER(category)::TAP_CATEGORY
    ELSE 'unknown'::TAP_CATEGORY
END;

ALTER TABLE tap
DROP COLUMN category;

ALTER TABLE tap
RENAME new_category TO category;

ALTER TABLE tap
ALTER COLUMN category DROP DEFAULT
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tap
ADD COLUMN old_category TEXT;

UPDATE tap
SET old_category = category::TEXT;

ALTER TABLE tap
DROP COLUMN category;

ALTER TABLE tap
RENAME COLUMN old_category TO category;

DROP TYPE TAP_CATEGORY;
-- +goose StatementEnd
