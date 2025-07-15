-- name: CreateTask :one
INSERT INTO tasks (
    task_day_id, user_id, title, description, duration
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetTask :one
SELECT * FROM tasks
WHERE id = $1 LIMIT 1;

-- name: ListTasks :many
SELECT * FROM tasks
ORDER BY title
LIMIT $1
OFFSET $2;

-- name: UpdateTask :exec
UPDATE tasks
  set title = $2,
  description = $3, 
  duration = $4, 
  completed = $5
WHERE id = $1;

-- name: DeleteTask :one
DELETE FROM tasks
WHERE id = $1
RETURNING *;