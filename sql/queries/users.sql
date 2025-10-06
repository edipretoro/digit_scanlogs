-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ?;

-- name: GetUserByUID :one
SELECT * FROM users WHERE uid = ?;

-- name: CreateUser :one
INSERT INTO users (id, uid, username, fullname)
VALUES (?, ?, ?, ?)
RETURNING *;