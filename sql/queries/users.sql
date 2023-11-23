-- name: CreateUser :one
INSERT INTO users (
        id,
        created_at,
        updated_at,
        username,
        fullname,
        hashed_password,
        bio,
        avatar,
        email,
        api_key
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        encode(sha256(random()::text::bytea), 'hex')
    )
RETURNING *;
-- name: UpdateProfile :one
UPDATE users
SET updated_at = NOW(),
    bio = COALESCE(sqlc.narg(bio), bio),
    username = COALESCE(sqlc.narg(username), username),
    avatar = COALESCE(sqlc.narg(avatar), avatar),
    fullname = COALESCE(sqlc.narg(fullname), fullname),
    email = COALESCE(sqlc.narg(email), email)
WHERE id = $1
RETURNING *;
-- name: UpdatePassword :one
UPDATE users
SET updated_at = NOW(),
    api_key = encode(sha256(random()::text::bytea), 'hex'),
    hashed_password = $2
WHERE id = $1
RETURNING *;
-- name: GetUser :one
SELECT *
FROM users
WHERE username = $1
    OR email = $2;
-- name: GetUserByAPIKey :one
SELECT *
FROM users
WHERE api_key = $1;
-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;