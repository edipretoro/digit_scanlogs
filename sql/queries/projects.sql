-- name: GetProjects :many
SELECT * 
FROM projects;

-- name: GetProject :one
SELECT * 
FROM projects 
WHERE id = ?;

-- name: GetProjectByPath :one
SELECT * 
FROM projects 
WHERE path = ?;

-- name: CreateProject :one
INSERT INTO projects (id, name, path, description, created_by)
VALUES (?, ?, ?, ?, ?)
RETURNING *;