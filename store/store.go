package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	sqlc "github.com/leandrohsilveira/simple-bank/database/sqlc"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	Queries *sqlc.Queries
	db *pgx.Conn
}

func NewStore(db *pgx.Conn) *Store {
	return &Store{
		db:      db,
		Queries: sqlc.New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) ExecTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := store.db.Begin(ctx)

	if err != nil {
		return err
	}

	queries := store.Queries.WithTx(tx)

	err = fn(queries)

	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit(ctx)
}