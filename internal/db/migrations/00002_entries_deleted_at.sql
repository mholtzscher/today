-- +goose Up
ALTER TABLE entries ADD COLUMN deleted_at TEXT NULL;
CREATE INDEX entries_deleted_at_idx ON entries(deleted_at);

-- +goose Down
DROP INDEX IF EXISTS entries_deleted_at_idx;
CREATE TABLE entries_rollback (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  text TEXT NOT NULL,
  created_at TEXT DEFAULT (datetime('now'))
);
INSERT INTO entries_rollback (id, text, created_at)
SELECT id, text, created_at
FROM entries;
DROP TABLE entries;
ALTER TABLE entries_rollback RENAME TO entries;
