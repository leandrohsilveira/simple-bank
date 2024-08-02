package database

import (
	"context"
	"os"
	"testing"

	"github.com/leandrohsilveira/simple-bank/configs"
	"github.com/leandrohsilveira/simple-bank/server/database"
)

var testQueries *database.Queries

func TestMain(m *testing.M) {
	ctx := context.Background()

	conn, teardown := SetupTestingDatabase(ctx, "../../database/migrations", configs.NewTestingDatabaseConfig())

	testQueries = database.New(conn)

	code := m.Run()

	teardown()

	os.Exit(code)
}
