// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: follow_habits.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createFollowHabit = `-- name: CreateFollowHabit :one
INSERT INTO follow_habits (id, created_at, updated_at, user_id, habit_id, promise_from, promise_end, processes, included_report) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, created_at, updated_at, user_id, habit_id, promise_from, promise_end, processes, included_report
`

type CreateFollowHabitParams struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	UserID         uuid.UUID
	HabitID        uuid.UUID
	PromiseFrom    time.Time
	PromiseEnd     time.Time
	Processes      []time.Time
	IncludedReport bool
}

func (q *Queries) CreateFollowHabit(ctx context.Context, arg CreateFollowHabitParams) (FollowHabit, error) {
	row := q.db.QueryRowContext(ctx, createFollowHabit,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.HabitID,
		arg.PromiseFrom,
		arg.PromiseEnd,
		pq.Array(arg.Processes),
		arg.IncludedReport,
	)
	var i FollowHabit
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.HabitID,
		&i.PromiseFrom,
		&i.PromiseEnd,
		pq.Array(&i.Processes),
		&i.IncludedReport,
	)
	return i, err
}

const deleteFollowHabit = `-- name: DeleteFollowHabit :exec
DELETE FROM follow_habits
WHERE habit_id = $1 AND user_id = $2
`

type DeleteFollowHabitParams struct {
	HabitID uuid.UUID
	UserID  uuid.UUID
}

func (q *Queries) DeleteFollowHabit(ctx context.Context, arg DeleteFollowHabitParams) error {
	_, err := q.db.ExecContext(ctx, deleteFollowHabit, arg.HabitID, arg.UserID)
	return err
}

const markExcute = `-- name: MarkExcute :one
UPDATE follow_habits
SET updated_at = NOW(),
    processes = array_append(processes, $3)
WHERE habit_id = $1 AND user_id = $2
RETURNING id, created_at, updated_at, user_id, habit_id, promise_from, promise_end, processes, included_report
`

type MarkExcuteParams struct {
	HabitID     uuid.UUID
	UserID      uuid.UUID
	ArrayAppend interface{}
}

func (q *Queries) MarkExcute(ctx context.Context, arg MarkExcuteParams) (FollowHabit, error) {
	row := q.db.QueryRowContext(ctx, markExcute, arg.HabitID, arg.UserID, arg.ArrayAppend)
	var i FollowHabit
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.HabitID,
		&i.PromiseFrom,
		&i.PromiseEnd,
		pq.Array(&i.Processes),
		&i.IncludedReport,
	)
	return i, err
}

const updateFollowHabit = `-- name: UpdateFollowHabit :one
UPDATE follow_habits
SET updated_at = NOW(),
    promise_from = COALESCE($3, promise_from),
    promise_end = COALESCE($4, promise_end),
    included_report = COALESCE($5, included_report)
WHERE habit_id = $1 AND user_id = $2
RETURNING id, created_at, updated_at, user_id, habit_id, promise_from, promise_end, processes, included_report
`

type UpdateFollowHabitParams struct {
	HabitID        uuid.UUID
	UserID         uuid.UUID
	PromiseFrom    sql.NullTime
	PromiseEnd     sql.NullTime
	IncludedReport sql.NullBool
}

func (q *Queries) UpdateFollowHabit(ctx context.Context, arg UpdateFollowHabitParams) (FollowHabit, error) {
	row := q.db.QueryRowContext(ctx, updateFollowHabit,
		arg.HabitID,
		arg.UserID,
		arg.PromiseFrom,
		arg.PromiseEnd,
		arg.IncludedReport,
	)
	var i FollowHabit
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.HabitID,
		&i.PromiseFrom,
		&i.PromiseEnd,
		pq.Array(&i.Processes),
		&i.IncludedReport,
	)
	return i, err
}