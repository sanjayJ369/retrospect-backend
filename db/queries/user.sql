-- name: CreateUser :one
INSERT INTO users (
  email, name
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: UpdateUser :exec
UPDATE users
  set name = $2,
  email = $3
WHERE id = $1;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING *;