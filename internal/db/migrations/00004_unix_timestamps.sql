-- +goose Up
-- Convert TEXT timestamps to INTEGER (Unix epoch)

PRAGMA foreign_keys=OFF;

CREATE TABLE entries_unix (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  text TEXT NOT NULL,
  created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
  archived_at INTEGER NULL
);

-- Convert existing TEXT timestamps to Unix epoch seconds
INSERT INTO entries_unix (id, text, created_at, archived_at)
SELECT 
  id,
  text,
  COALESCE(
    CAST(strftime('%s', created_at) AS INTEGER),
    strftime('%s', 'now')
  ),
  CASE 
    WHEN archived_at IS NOT NULL THEN CAST(strftime('%s', archived_at) AS INTEGER)
    ELSE NULL
  END
FROM entries;

DROP INDEX IF EXISTS entries_archived_at_idx;
DROP TABLE entries;
ALTER TABLE entries_unix RENAME TO entries;
CREATE INDEX entries_archived_at_idx ON entries(archived_at);

PRAGMA foreign_keys=ON;

-- +goose Down
-- Convert INTEGER timestamps back to TEXT

PRAGMA foreign_keys=OFF;

CREATE TABLE entries_text (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  text TEXT NOT NULL,
  created_at TEXT DEFAULT (datetime('now')),
  archived_at TEXT NULL
);

INSERT INTO entries_text (id, text, created_at, archived_at)
SELECT 
  id,
  text,
  datetime(created_at, 'unixepoch'),
  CASE 
    WHEN archived_at IS NOT NULL THEN datetime(archived_at, 'unixepoch')
    ELSE NULL
  END
FROM entries;

DROP INDEX IF EXISTS entries_archived_at_idx;
DROP TABLE entries;
ALTER TABLE entries_text RENAME TO entries;
CREATE INDEX entries_archived_at_idx ON entries(archived_at);

PRAGMA foreign_keys=ON;
