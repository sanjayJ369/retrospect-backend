-- name: CreateChallenge :one
INSERT INTO challenges (
  title, user_id, description, end_date, duration
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetChallenge :one
SELECT * FROM challenges
WHERE id = $1 LIMIT 1;

-- name: ListChallenges :many
SELECT * FROM challenges
ORDER BY start_date
LIMIT $1
OFFSET $2;

-- name: UpdateChallenge :exec
UPDATE challenges
  set title = $2,
  description = $3, 
  end_date = $4, 
  duration = $5, 
  active = $6
WHERE id = $1;

-- name: DeleteChallenge :one
DELETE FROM users
WHERE id = $1
RETURNING *;