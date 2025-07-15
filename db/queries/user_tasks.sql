-- name: ListTaskDaysByUserId :many
SELECT * FROM task_days
WHERE user_id = $1
ORDER BY date
LIMIT $2
OFFSET $3;

-- name: ListTasksByTaskDayId :many
SELECT * FROM tasks
WHERE task_day_id = $1;