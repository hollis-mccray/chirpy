-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: UserByEmail :one
SELECT * from users
WHERE email = $1 limit 1;

-- name: GetUserFromRefreshToken :one
SELECT users.* from users
INNER JOIN refresh_tokens
ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1
LIMIT 1;

-- name: UpdatePassword :one
UPDATE users
SET email = $2, hashed_password = $3
WHERE id = $1
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: UpdateToRed :one
UPDATE users
SET is_chirpy_red = true
where id = $1
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: ResetUsers :exec
DELETE FROM users;
