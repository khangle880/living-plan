// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: projects.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createProject = `-- name: CreateProject :one
INSERT INTO projects (id, created_at, updated_at, title, creator, purpose, description, background)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, created_at, updated_at, title, creator, purpose, description, background
`

type CreateProjectParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Creator     uuid.UUID
	Purpose     sql.NullString
	Description sql.NullString
	Background  sql.NullString
}

func (q *Queries) CreateProject(ctx context.Context, arg CreateProjectParams) (Project, error) {
	row := q.db.QueryRowContext(ctx, createProject,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Title,
		arg.Creator,
		arg.Purpose,
		arg.Description,
		arg.Background,
	)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Creator,
		&i.Purpose,
		&i.Description,
		&i.Background,
	)
	return i, err
}

const deleteProject = `-- name: DeleteProject :exec
DELETE FROM projects WHERE id = $1 AND creator = $2
`

type DeleteProjectParams struct {
	ID      uuid.UUID
	Creator uuid.UUID
}

func (q *Queries) DeleteProject(ctx context.Context, arg DeleteProjectParams) error {
	_, err := q.db.ExecContext(ctx, deleteProject, arg.ID, arg.Creator)
	return err
}

const getProjects = `-- name: GetProjects :many
SELECT id, created_at, updated_at, title, creator, purpose, description, background FROM projects 
OFFSET $1
LIMIT $2
`

type GetProjectsParams struct {
	Offset int32
	Limit  int32
}

func (q *Queries) GetProjects(ctx context.Context, arg GetProjectsParams) ([]Project, error) {
	rows, err := q.db.QueryContext(ctx, getProjects, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Project
	for rows.Next() {
		var i Project
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Creator,
			&i.Purpose,
			&i.Description,
			&i.Background,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProjectsByUser = `-- name: GetProjectsByUser :many
SELECT projects.id, projects.created_at, projects.updated_at, projects.title, projects.creator, projects.purpose, projects.description, projects.background FROM projects
JOIN follow_projects ON projects.id = follow_projects.project_id
WHERE follow_projects.user_id = $1
OFFSET $2
LIMIT $3
`

type GetProjectsByUserParams struct {
	UserID uuid.UUID
	Offset int32
	Limit  int32
}

func (q *Queries) GetProjectsByUser(ctx context.Context, arg GetProjectsByUserParams) ([]Project, error) {
	rows, err := q.db.QueryContext(ctx, getProjectsByUser, arg.UserID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Project
	for rows.Next() {
		var i Project
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Creator,
			&i.Purpose,
			&i.Description,
			&i.Background,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateProject = `-- name: UpdateProject :one
UPDATE projects
SET updated_at = NOW(),
    title = COALESCE($3, title),
    purpose = COALESCE($4, purpose),
    description = COALESCE($5, description),
    background = COALESCE($6, background)
WHERE id = $1 AND creator = $2
RETURNING id, created_at, updated_at, title, creator, purpose, description, background
`

type UpdateProjectParams struct {
	ID          uuid.UUID
	Creator     uuid.UUID
	Title       sql.NullString
	Purpose     sql.NullString
	Description sql.NullString
	Background  sql.NullString
}

func (q *Queries) UpdateProject(ctx context.Context, arg UpdateProjectParams) (Project, error) {
	row := q.db.QueryRowContext(ctx, updateProject,
		arg.ID,
		arg.Creator,
		arg.Title,
		arg.Purpose,
		arg.Description,
		arg.Background,
	)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Creator,
		&i.Purpose,
		&i.Description,
		&i.Background,
	)
	return i, err
}
