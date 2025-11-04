-- name: SetRefToken :one
INSERT INTO refresh_tokens (
    token,
    user_id,
    expires_at,
    revoked_at
) VALUES (
    $1,
    $2,
    $3,
    NULL
) RETURNING *;

-- name: GetRefToken :one
SELECT * FROM refresh_tokens WHERE token = $1;

-- name: GetUserFromRefToken :one
SELECT user_id FROM refresh_tokens WHERE token = $1;

-- name: RevokeRefToken :one
UPDATE refresh_tokens SET revoked_at = NOW(), updated_at = NOW() WHERE token = $1 RETURNING *;