-- name: CreateEntry :exec
INSERT INTO entries (text)
VALUES (?);

-- name: CreateEntryAt :exec
INSERT INTO entries (text, created_at)
VALUES (?, ?);

-- name: GetEntry :one
SELECT id, text, created_at, archived_at
FROM entries
WHERE id = ?
LIMIT 1;

-- name: ListEntriesSince :many
SELECT id, text, created_at, archived_at
FROM entries
WHERE datetime(created_at, 'unixepoch') >= datetime('now', CAST(? AS TEXT))
AND archived_at IS NULL
ORDER BY created_at DESC, id DESC;

-- name: ListEntriesSinceAll :many
SELECT id, text, created_at, archived_at
FROM entries
WHERE datetime(created_at, 'unixepoch') >= datetime('now', CAST(? AS TEXT))
ORDER BY created_at DESC, id DESC;

-- name: ArchiveEntry :execrows
UPDATE entries
SET archived_at = strftime('%s', 'now')
WHERE id = ? AND archived_at IS NULL;

-- name: RestoreEntry :execrows
UPDATE entries
SET archived_at = NULL
WHERE id = ? AND archived_at IS NOT NULL;
