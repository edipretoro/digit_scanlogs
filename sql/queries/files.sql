-- name: GetFiles :many
SELECT * FROM files;

-- name: GetFile :one
SELECT * FROM files WHERE id = ?;

-- name: GetFileByPath :one
SELECT * FROM files WHERE path = ?;

-- name: CreateFile :one
INSERT INTO files (id, project_id, user_id, name, path, size, mode, modtime, sha512, description)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;