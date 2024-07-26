package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leandrohsilveira/simple-bank/configs"
)

func SetupTestingDatabase(ctx context.Context, migrationsSource string, config configs.DatabaseConfig) (*pgxpool.Pool, func()) {
	dbSource := config.GetDbSource()

	dbConfig, err := pgxpool.ParseConfig(dbSource)

	dbConfig.MaxConns = int32(config.MaxConns)

	if err != nil {
		log.Fatalln("cannot parse db config:", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)

	if err != nil {
		log.Fatalln("cannot create db connection pool:", err)
	}

	migrator, err := NewMigrate(migrationsSource, dbSource)

	if err != nil {
		pool.Close()
		log.Fatalln("cannot create migration instance:", err)
	}

	err = migrator.Down()

	if err != nil {
		pool.Close()
		log.Println("cannot run migration down:", err.Error())
	}

	err = migrator.Up()

	if err != nil {
		migrator.Close()
		pool.Close()
		log.Fatalln("cannot run migration up:", err)
	}

	teardown := func() {
		pool.Close()
	
		if err != nil {
			migrator.Close()
			log.Fatalln("cannot close db connection:", err)
		}
	
		sourceErr, err := migrator.Close()
	
		if err != nil {
			log.Fatalln("cannot close migration instance:", err)
		}
	
		if sourceErr != nil {
			log.Fatalln("cannot close migration source:", sourceErr)
		}
	}

	return pool, teardown
}