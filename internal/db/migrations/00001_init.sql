-- +goose Up
CREATE TABLE entries (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  text TEXT NOT NULL,
  created_at TEXT DEFAULT (datetime('now'))
);

-- +goose Down
DROP TABLE entries;