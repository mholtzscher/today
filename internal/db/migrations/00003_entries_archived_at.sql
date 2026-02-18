-- +goose Up
CREATE TABLE entries_archived (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  text TEXT NOT NULL,
  created_at TEXT DEFAULT (datetime('now')),
  archived_at TEXT NULL
);
INSERT INTO entries_archived (id, text, created_at, archived_at)
SELECT id, text, created_at, deleted_at
FROM entries;
DROP INDEX IF EXISTS entries_deleted_at_idx;
DROP TABLE entries;
ALTER TABLE entries_archived RENAME TO entries;
CREATE INDEX entries_archived_at_idx ON entries(archived_at);

-- +goose Down
CREATE TABLE entries_deleted (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  text TEXT NOT NULL,
  created_at TEXT DEFAULT (datetime('now')),
  deleted_at TEXT NULL
);
INSERT INTO entries_deleted (id, text, created_at, deleted_at)
SELECT id, text, created_at, archived_at
FROM entries;
DROP INDEX IF EXISTS entries_archived_at_idx;
DROP TABLE entries;
ALTER TABLE entries_deleted RENAME TO entries;
CREATE INDEX entries_deleted_at_idx ON entries(deleted_at);
