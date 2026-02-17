-- +goose Up
-- +goose StatementBegin
CREATE TABLE reply (
  id SERIAL PRIMARY KEY,
  message_id INTEGER NOT NULL REFERENCES message (id) ON DELETE CASCADE,
  name TEXT,
  message TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE reply;
-- +goose StatementEnd
