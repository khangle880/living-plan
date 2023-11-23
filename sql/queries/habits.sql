-- name: CreateHabit :one
INSERT INTO habits (
        id,
        created_at,
        updated_at,
        title,
        creator,
        purpose,
        description,
        icon,
        background,
        images,
        urls,
        time_in_day,
        loop_week,
        loop_month,
        recommend_duration
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING *;

-- name: UpdateHabit :one
UPDATE habits
SET updated_at = NOW(),
    title = COALESCE(sqlc.narg(title), title),
    purpose = COALESCE(sqlc.narg(purpose), purpose),
    description = COALESCE(sqlc.narg(description), description),
    icon = COALESCE(sqlc.narg(icon), icon),
    background = COALESCE(sqlc.narg(background), background),
    images = COALESCE(sqlc.narg(images), images),
    urls = COALESCE(sqlc.narg(urls), urls),
    time_in_day = COALESCE(sqlc.narg(time_in_day), time_in_day),
    loop_week = COALESCE(sqlc.narg(loop_week), loop_week),
    loop_month = COALESCE(sqlc.narg(loop_month), loop_month),
    recommend_duration = COALESCE(sqlc.narg(recommend_duration), recommend_duration)
WHERE id = $1 AND creator = $2
RETURNING *;

-- name: GetHabits :many
SELECT *
FROM habits OFFSET $1
LIMIT $2;
-- name: GetHabitsByUser :many
SELECT habits.*
FROM habits
    JOIN follow_habits ON habits.id = follow_habits.habit_id
WHERE follow_habits.user_id = $1 OFFSET $2
LIMIT $3;
-- name: DeleteHabit :exec
DELETE FROM habits
WHERE id = $1 AND creator = $2;