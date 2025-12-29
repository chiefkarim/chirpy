-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES ($1, now(), now(), $2, $3)
RETURNING *;

-- name: GetRefreshToken :one
SELECT
    token,
    created_at,
    updated_at,
    user_id,
    expires_at,
    revoked_at
FROM refresh_tokens
WHERE token = $1;

-- name: ExpireToken :one
UPDATE refresh_tokens
SET
    revoked_at = $1,
    updated_at = now()
WHERE token = $2
RETURNING *;
