package main

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leandrohsilveira/simple-bank/configs"

	fiberRecover "github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {

	dbConfig := configs.NewDatabaseConfig()

	poolConfig, err := pgxpool.ParseConfig(dbConfig.GetDbSource())

	if err != nil {
		log.Fatal("cannot parse db config:", err)
	}

	poolConfig.MaxConns = int32(dbConfig.MaxConns)

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)

	if err != nil {
		log.Fatal("cannot create db connection pool:", err)
	}

	app := fiber.New(fiber.Config{
		AppName:       "Simple Bank",
		CaseSensitive: true,
	})

	app.Use(fiberRecover.New(fiberRecover.Config{
		EnableStackTrace: true,
	}))
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${locals:requestid} ${ip} ${status} - ${latency} ${method} ${path} ${error}\n",
	}))

	app.Get("/*", static.New("./client/dist"))
	app.Get("/*", static.New("./client/static"))

	app.Use(AppRouter(pool))

	log.Fatal(app.Listen(":3000"))
}
