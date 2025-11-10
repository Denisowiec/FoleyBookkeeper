-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    hashed_password
) VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByName :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetAllUsers :many
SELECT id, created_at, updated_at, username, email FROM users;

-- name: DeleteUser :one
DELETE FROM users WHERE id = $1 RETURNING *;

-- name: UpdateUser :one
UPDATE users SET
    updated_at = NOW(),
    username = $2,
    email = $3
WHERE id = $1 RETURNING *;

-- name: ChangePassword :one
UPDATE users SET hashed_password = $2 WHERE id = $1 RETURNING *;
