package database

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)


func NewMigrate(relativePath, dbSource string) (*migrate.Migrate, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return nil, fmt.Errorf("cannot get current working directory: %w", err)
	}

	appPath, err := filepath.Abs(path.Join(cwd, relativePath))

	if err != nil {
		return nil, fmt.Errorf("cannot get absolute path: %w", err)
	}

	return migrate.New(path.Join("file://", appPath), dbSource)
}