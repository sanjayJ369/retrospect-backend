-- name: CreateUser :one
INSERT INTO users (
  email, name, hashed_password
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByName :one 
SELECT * FROM users 
WHERE name = $1 LIMIT 1; 

-- name: GetUserByEmail :one 
SELECT * FROM users
WHERE email = $1 LIMIT 1;

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

-- name: UpdateUserIsVerified :one
UPDATE users
SET
  is_verified = $2,
  updated_at = NOW()
WHERE
  id = $1
RETURNING *;