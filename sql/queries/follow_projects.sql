-- name: CreateFollowProject :one
INSERT INTO follow_projects (id, created_at, updated_at, user_id, project_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
-- name: DeleteFollowProject :exec
DELETE FROM follow_projects
WHERE project_id = $1
    AND user_id = $2;