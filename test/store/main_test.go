package store

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leandrohsilveira/simple-bank/configs"
	"github.com/leandrohsilveira/simple-bank/server/store"
	"github.com/leandrohsilveira/simple-bank/test/database"
)

var testDbPool *pgxpool.Pool

func CreateStore(ctx context.Context) (*store.Store, func(), error) {
	conn, err := testDbPool.Acquire(ctx)

	if err != nil {
		return nil, nil, err
	}

	return store.NewStore(conn.Conn()), conn.Release, nil
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	var teardown func()
	testDbPool, teardown = database.SetupTestingDatabase(ctx, "../../database/migrations", configs.NewTestingDatabaseConfig())

	code := m.Run()

	teardown()

	os.Exit(code)
}
