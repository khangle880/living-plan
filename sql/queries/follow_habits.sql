-- name: CreateFollowHabit :one
INSERT INTO follow_habits (id, created_at, updated_at, user_id, habit_id, promise_from, promise_end, processes, included_report) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: MarkExcute :one
UPDATE follow_habits
SET updated_at = NOW(),
    processes = array_append(processes, $3)
WHERE habit_id = $1 AND user_id = $2
RETURNING *;

-- name: UpdateFollowHabit :one
UPDATE follow_habits
SET updated_at = NOW(),
    promise_from = COALESCE(sqlc.narg(promise_from), promise_from),
    promise_end = COALESCE(sqlc.narg(promise_end), promise_end),
    included_report = COALESCE(sqlc.narg(included_report), included_report)
WHERE habit_id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteFollowHabit :exec
DELETE FROM follow_habits
WHERE habit_id = $1 AND user_id = $2;