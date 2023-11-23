-- name: CreateTask :one
INSERT INTO tasks (
        id,
        created_at,
        updated_at,
        datetime_exc,
        title,
        purpose,
        description,
        images,
        urls,
        status,
        user_id,
        project_id,
        parent_id
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: UpdateTask :one
UPDATE tasks
SET updated_at = NOW(),
    datetime_exc = COALESCE(sqlc.narg(datetime_exc), datetime_exc),
    title = COALESCE(sqlc.narg(title), title),
    purpose = COALESCE(sqlc.narg(purpose), purpose),
    description = COALESCE(sqlc.narg(description), description),
    images = COALESCE(sqlc.narg(images), images),
    urls = COALESCE(sqlc.narg(urls), urls),
    project_id = COALESCE(sqlc.narg(project_id), project_id),
    parent_id = COALESCE(sqlc.narg(parent_id), parent_id)
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: UpdateStatus :one
UPDATE tasks
SET status = $2
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: GetTasksByUser :many
SELECT * FROM tasks 
WHERE user_id = $1 OFFSET $2
LIMIT $3;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1 AND user_id = $2;
