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

-- name: UpdateUserName :one
UPDATE users
SET
  name = $2,
  updated_at = NOW()
WHERE
  id = $1
RETURNING *;

-- name: UpdateUserEmail :one
UPDATE users
SET
  email = $2,
  updated_at = NOW()
WHERE
  id = $1
RETURNING *;

-- name: UpdateUserTimezone :one
UPDATE users
SET
  timezone = $2,
  updated_at = NOW()
WHERE
  id = $1
RETURNING *;


-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING *;