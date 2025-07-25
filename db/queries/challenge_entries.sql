-- name: CreateChallengeEntry :one
INSERT INTO challenge_entries (
  challenge_id
) VALUES (
  $1
)
RETURNING *;

-- name: GetChallengeEntry :one
SELECT * FROM challenge_entries
WHERE id = $1 LIMIT 1;

-- name: ListChallengeEntries :many
SELECT * FROM challenge_entries
ORDER BY created_at
LIMIT $1
OFFSET $2;

-- name: UpdateChallengeEntry :one
UPDATE challenge_entries
  set completed = $2
WHERE id = $1
RETURNING *;;

-- name: DeleteChallengeEntry :one
DELETE FROM challenge_entries
WHERE id = $1
RETURNING *;