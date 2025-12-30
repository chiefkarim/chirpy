-- name: CreateChirp :one
INSERT INTO chirps (id, body, created_at, updated_at, user_id)
VALUES (gen_random_uuid(), $1, now(), now(), $2)
RETURNING *;

-- name: GetAllChirps :many
SELECT
    id,
    body,
    created_at,
    updated_at,
    user_id
FROM chirps
ORDER BY created_at;

-- name: GetChirp :one
SELECT
    id,
    body,
    created_at,
    updated_at,
    user_id
FROM chirps
WHERE id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE user_id = $1;
