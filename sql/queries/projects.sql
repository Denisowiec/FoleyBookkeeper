-- name: CreateProject :one
INSERT INTO projects (
    title,
    client_id
) VALUES (
    $1,
    $2
) RETURNING *;

-- name: GetProjectByID :one
SELECT * FROM projects WHERE id=$1;

-- name: GetProjectByTitle :one
SELECT * FROM projects WHERE title=$1;

-- name: GetAllProjects :many
SELECT * FROM projects;

-- name: GetProjectsByClient :many
SELECT * FROM projects WHERE client_id=$1;

-- name: UpdateProject :one
UPDATE projects SET
    updated_at = NOW(),
    title = $2,
    client_id = $3
WHERE id = $1 RETURNING *;

-- name: DeleteProject :one
DELETE FROM projects WHERE id=$1 RETURNING *;
