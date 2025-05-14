-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;


-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: ChangeUserPassword :one
UPDATE users
SET 
    hashed_password = $1,
    email = $2
WHERE id = $3
RETURNING *;


-- name: UpdateUser :one
UPDATE users
SET    
    created_at = $1,
    updated_at = $2,
    hashed_password = $3,
    email = $4
WHERE id = $5
RETURNING *;


-- name: UpgradeUser :exec
UPDATE users
SET 
    is_chirpy_red = true
WHERE id = $1
RETURNING *;