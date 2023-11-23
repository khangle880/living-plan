-- name: CreateProject :one
INSERT INTO projects (id, created_at, updated_at, title, creator, purpose, description, background)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateProject :one
UPDATE projects
SET updated_at = NOW(),
    title = COALESCE(sqlc.narg(title), title),
    purpose = COALESCE(sqlc.narg(purpose), purpose),
    description = COALESCE(sqlc.narg(description), description),
    background = COALESCE(sqlc.narg(background), background)
WHERE id = $1 AND creator = $2
RETURNING *;

-- name: GetProjectsByUser :many
SELECT projects.* FROM projects
JOIN follow_projects ON projects.id = follow_projects.project_id
WHERE follow_projects.user_id = $1
OFFSET $2
LIMIT $3;

-- name: GetProjects :many
SELECT * FROM projects 
OFFSET $1
LIMIT $2;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = $1 AND creator = $2;