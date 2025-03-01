-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS season (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    start TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    "end" TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS scan (
  id SERIAL PRIMARY KEY,
  scan_time TIMESTAMP WITH TIME ZONE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS scan;

DROP TABLE IF EXISTS season;
-- +goose StatementEnd
