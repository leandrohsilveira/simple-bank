package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	user_domain "github.com/leandrohsilveira/simple-bank/server/domain/user"
	"github.com/leandrohsilveira/simple-bank/server/store"
)

func AppRouter(pool *pgxpool.Pool) fiber.Router {

	app := fiber.New()

	app.Use("/users", store.WithStore(pool, user_domain.UserRouter()))

	return app
}
