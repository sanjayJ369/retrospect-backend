// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: sessions.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createSession = `-- name: CreateSession :one
INSERT INTO sessions (
  id,
  user_id  ,
  refresh_token ,
  user_agent, 
  client_ip , 
  is_blocked , 
  expires_at 
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, user_id, refresh_token, user_agent, client_ip, is_blocked, created_at, expires_at
`

type CreateSessionParams struct {
	ID           pgtype.UUID      `json:"id"`
	UserID       pgtype.UUID      `json:"user_id"`
	RefreshToken string           `json:"refresh_token"`
	UserAgent    string           `json:"user_agent"`
	ClientIp     string           `json:"client_ip"`
	IsBlocked    bool             `json:"is_blocked"`
	ExpiresAt    pgtype.Timestamp `json:"expires_at"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRow(ctx, createSession,
		arg.ID,
		arg.UserID,
		arg.RefreshToken,
		arg.UserAgent,
		arg.ClientIp,
		arg.IsBlocked,
		arg.ExpiresAt,
	)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const getSessions = `-- name: GetSessions :one
SELECT id, user_id, refresh_token, user_agent, client_ip, is_blocked, created_at, expires_at FROM sessions 
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetSessions(ctx context.Context, id pgtype.UUID) (Session, error) {
	row := q.db.QueryRow(ctx, getSessions, id)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}
