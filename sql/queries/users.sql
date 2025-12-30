-- name: CreateUser :one
INSERT INTO users (id, email, hashed_password, created_at, updated_at)
VALUES (gen_random_uuid(), $1, $2, now(), now())
RETURNING *;

-- name: GetUserByEmail :one
SELECT
    id,
    email,
    created_at,
    updated_at,
    hashed_password,
    is_chirpy_red
FROM users
WHERE email = $1;

-- name: ChangeUserDetails :one
UPDATE users
SET
    email = $1,
    hashed_password = $2,
    updated_at = now()
WHERE id = $3
RETURNING *;

-- name: UpgradeUser :one
UPDATE users
SET
    is_chirpy_red = true
WHERE id = $1
RETURNING *;
