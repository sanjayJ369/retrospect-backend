// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: user_challenges.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const listChallengeEntriesByChallengeId = `-- name: ListChallengeEntriesByChallengeId :many
SELECT id, challenge_id, date, completed, created_at FROM challenge_entries
WHERE challenge_id = $1
ORDER BY challenge_entries.date
LIMIT $2
OFFSET $3
`

type ListChallengeEntriesByChallengeIdParams struct {
	ChallengeID pgtype.UUID `json:"challenge_id"`
	Limit       int32       `json:"limit"`
	Offset      int32       `json:"offset"`
}

func (q *Queries) ListChallengeEntriesByChallengeId(ctx context.Context, arg ListChallengeEntriesByChallengeIdParams) ([]ChallengeEntry, error) {
	rows, err := q.db.Query(ctx, listChallengeEntriesByChallengeId, arg.ChallengeID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ChallengeEntry{}
	for rows.Next() {
		var i ChallengeEntry
		if err := rows.Scan(
			&i.ID,
			&i.ChallengeID,
			&i.Date,
			&i.Completed,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listChallengesByUser = `-- name: ListChallengesByUser :many
SELECT id, title, user_id, description, start_date, end_date, active, created_at FROM challenges
WHERE user_id = $1
ORDER BY start_date
LIMIT $2
OFFSET $3
`

type ListChallengesByUserParams struct {
	UserID pgtype.UUID `json:"user_id"`
	Limit  int32       `json:"limit"`
	Offset int32       `json:"offset"`
}

func (q *Queries) ListChallengesByUser(ctx context.Context, arg ListChallengesByUserParams) ([]Challenge, error) {
	rows, err := q.db.Query(ctx, listChallengesByUser, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Challenge{}
	for rows.Next() {
		var i Challenge
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.UserID,
			&i.Description,
			&i.StartDate,
			&i.EndDate,
			&i.Active,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
