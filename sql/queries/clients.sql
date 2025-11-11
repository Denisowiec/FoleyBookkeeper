-- name: CreateClient :one
INSERT INTO clients (
    client_name,
    email,
    notes
) VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: GetClientByID :one
SELECT * FROM clients WHERE id=$1;

-- name: GetClientByName :one
SELECT * FROM clients WHERE client_name=$1;

-- name: GetAllClients :many
SELECT * FROM clients;

-- name: UpdateClient :one
UPDATE clients SET
    updated_at=NOW(),
    client_name=$2,
    email=$3,
    notes=$4
WHERE id=$1 RETURNING *;

-- name: DeleteClient :one
DELETE FROM clients WHERE id=$1 RETURNING *;
