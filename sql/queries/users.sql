-- name: CreateUser :one
INSERT INTO users (id, email, hashed_password, created_at, updated_at)
VALUES (gen_random_uuid(), $1, $2, now(), now())
RETURNING *;

-- name: GetUserByEmail :one
SELECT id, email, hashed_password, created_at, updated_at FROM users WHERE email = $1;
