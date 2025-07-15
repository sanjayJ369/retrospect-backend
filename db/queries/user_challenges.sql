-- name: ListChallengesByUser :many
SELECT * FROM challenges
WHERE user_id = $1
ORDER BY start_date
LIMIT $2
OFFSET $3;

-- name: ListChallengeEntriesByChallengeId :many
SELECT * FROM challenge_entries
WHERE challenge_id = $1
ORDER BY challenge_entries.date
LIMIT $2
OFFSET $3;

