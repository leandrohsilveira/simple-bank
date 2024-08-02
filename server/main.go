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
	"github.com/leandrohsilveira/simple-bank/server/view"
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
		Format: "[${time}] ${latency} - ip:${ip} accept:${accept} ${method} ${path} in:${bytesReceived}b - ${status} content-type:${respHeader:content-type} out:${bytesSent}b: ${error}\n",
		CustomTags: map[string]logger.LogFunc{
			"accept": func(output logger.Buffer, c fiber.Ctx, data *logger.Data, _ string) (int, error) {
				return output.WriteString(view.Accept(c))
			},
		},
	}))

	app.Get("/*", static.New("./client/dist"))
	app.Get("/*", static.New("./client/static"))

	app.Use(AppRouter(pool))

	log.Fatal(app.Listen(":3000"))
}
