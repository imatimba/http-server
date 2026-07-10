-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password_hash)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: LookupUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: LookupUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users
WHERE TRUE;

-- name: UpdateUser :one
UPDATE users
SET email = $2, password_hash = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpgradeUserToChirpyRed :exec
UPDATE users
SET is_chirpy_red = TRUE, updated_at = NOW()
WHERE id = $1;