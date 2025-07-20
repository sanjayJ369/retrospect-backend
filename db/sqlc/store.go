package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Store interface {
	Querier
}

type SQLStore struct {
	*Queries
	conn *pgx.Conn
}

func NewStore(conn *pgx.Conn) Store {
	return &SQLStore{
		Queries: New(conn),
		conn:    conn,
	}
}

func (store *SQLStore) execTxn(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("txErr: %w, rbErr: %w", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
