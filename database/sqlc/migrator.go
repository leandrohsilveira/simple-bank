package database

import (
	"log"
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
		log.Fatalln("cannot get current working directory:", err)
	}

	appPath, err := filepath.Abs(path.Join(cwd, relativePath))

	if err != nil {
		log.Fatalln("cannot get absolute path:", err)
	}

	return migrate.New(path.Join("file://", appPath), dbSource)
}