-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS season (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    start TIMESTAMP NOT NULL,
    end TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS scan (
  id INTEGER PRIMARY KEY,
  scan_time TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS scan;

DROP TABLE IF EXISTS season;
-- +goose StatementEnd
