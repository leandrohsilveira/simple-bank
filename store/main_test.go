package store

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leandrohsilveira/simple-bank/configs"
	database "github.com/leandrohsilveira/simple-bank/database/sqlc"
)

var testDbPool *pgxpool.Pool

func CreateStore(ctx context.Context) (*Store, func(), error) {
	conn, err := testDbPool.Acquire(ctx)

	if err != nil {
		return nil, nil, err
	}

	return NewStore(conn.Conn()), conn.Release, nil
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	var teardown func()
	testDbPool, teardown = database.SetupTestingDatabase(ctx, "../database/migrations", configs.NewTestingDatabaseConfig())

	code := m.Run()

	teardown()

	os.Exit(code)
}