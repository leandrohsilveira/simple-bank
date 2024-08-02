package store

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	sqlc "github.com/leandrohsilveira/simple-bank/server/database"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	Queries *sqlc.Queries
	db      *pgx.Conn
}

var StoreCtxKey *Store = &Store{}

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

func WithStore(pool *pgxpool.Pool, router *fiber.App) fiber.Router {
	app := fiber.New()

	app.Use(
		func(c fiber.Ctx) error {
			conn, err := pool.Acquire(c.UserContext())

			if err != nil {
				return err
			}

			defer conn.Release()

			store := NewStore(conn.Conn())

			c.Locals(StoreCtxKey, store)

			return c.Next()
		},
	)

	app.Use(router)

	return app
}
