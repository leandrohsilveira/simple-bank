package database

import (
	"context"
	"os"
	"testing"

	"github.com/leandrohsilveira/simple-bank/configs"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	ctx := context.Background()

	conn, teardown := SetupTestingDatabase(ctx, "../migrations", configs.NewTestingDatabaseConfig())

	testQueries = New(conn)

	code := m.Run()

	teardown()

	os.Exit(code)
}