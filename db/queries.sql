-- name: CreateUser :execresult
INSERT INTO users (email)
VALUES (?);

-- name: GetUser :one
SELECT id, email
FROM users
WHERE id = ?
LIMIT 1;

-- name: ListUsers :many
SELECT id, email
FROM users
ORDER BY id;