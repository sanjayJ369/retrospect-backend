-- name: CreateTaskDay :one
INSERT INTO task_days (
    user_id
) VALUES (
  $1
)
RETURNING *;

-- name: GetTaskDay :one
SELECT * FROM task_days
WHERE id = $1 LIMIT 1;

-- name: ListTaskDays :many
SELECT * FROM task_days
ORDER BY task_days.date
LIMIT $1
OFFSET $2;

-- name: DeleteTaskDay :one
DELETE FROM task_days
WHERE id = $1
RETURNING *;