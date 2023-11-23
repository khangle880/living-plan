-- name: CreateDiary :one
INSERT INTO diaries (
        id,
        created_at,
        updated_at,
        datetime_exc,
        user_id,
        title,
        description,
        images,
        urls,
        icon,
        background
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;
-- name: GetDiariesByUser :many
SELECT *
FROM diaries
WHERE user_id = $1 OFFSET $2
LIMIT $3;
-- name: UpdateDiary :one
UPDATE diaries
SET updated_at = NOW(),
    datetime_exc = COALESCE(sqlc.narg(datetime_exc), datetime_exc),
    title = COALESCE(sqlc.narg(title), title),
    description = COALESCE(sqlc.narg(description), description),
    images = COALESCE(sqlc.narg(images), images),
    urls = COALESCE(sqlc.narg(urls), urls),
    icon = COALESCE(sqlc.narg(icon), icon),
    background = COALESCE(sqlc.narg(background), background)
WHERE id = $1 AND user_id = $2
RETURNING *;
-- name: DeleteDiary :exec
DELETE FROM diaries
WHERE id = $1
    AND user_id = $2;