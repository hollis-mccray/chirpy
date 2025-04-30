-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetRefreshTokenFromToken :one
SELECT * from refresh_tokens
WHERE token = $1
LIMIT 1;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = $2, updated_at = $2
where token = $1;