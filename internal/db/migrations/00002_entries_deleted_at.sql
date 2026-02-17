-- +goose Up
ALTER TABLE entries ADD COLUMN deleted_at TEXT NULL;
CREATE INDEX entries_deleted_at_idx ON entries(deleted_at);

-- +goose Down
DROP INDEX IF EXISTS entries_deleted_at_idx;
