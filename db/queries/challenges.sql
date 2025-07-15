-- name: CreateChallenge :one
INSERT INTO challenges (
  title, user_id, description, end_date
) VALUES (
  $1, $2, $3, $4
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

-- name: UpdateChallengeTitle :one
UPDATE challenges
SET
  title = $2
WHERE
  id = $1
RETURNING *;

-- name: UpdateChallengeDescription :one
UPDATE challenges
SET
  description = $2
WHERE
  id = $1
RETURNING *;

-- name: UpdateChallengeEndDate :one
UPDATE challenges
SET
  end_date = $2
WHERE
  id = $1
RETURNING *;

-- name: UpdateChallengeActiveStatus :one
UPDATE challenges
SET
  active = $2
WHERE
  id = $1
RETURNING *;

-- name: UpdateChallengeDetails :one
-- This query updates multiple common fields together
UPDATE challenges
SET
  title = $2,
  description = $3,
  end_date = $4
WHERE
  id = $1
RETURNING *;


-- name: DeleteChallenge :one
DELETE FROM challenges
WHERE id = $1
RETURNING *;