-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tap (
    id INTEGER PRIMARY KEY,
    order_id INTEGER NOT NULL,
    order_created_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tap;
-- +goose StatementEnd
