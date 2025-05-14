-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (created_at, updated_at, user_id, refresh_token, expires_at, revoked_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5, 
    $6
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens WHERE refresh_token = $1 LIMIT 1;

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET
    revoked_at = $1,
    updated_at = $2
WHERE refresh_token = $3
RETURNING *;