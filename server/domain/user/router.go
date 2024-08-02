package user_domain

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
	"github.com/leandrohsilveira/simple-bank/server/view"
)

func UserRouter() *fiber.App {
	app := fiber.New()

	app.Get("/",
		view.Json(func(ctx fiber.Ctx) error {
			userService, err := NewUserService(ctx)

			if err != nil {
				return err
			}

			users, err := userService.UserTableData(ctx.UserContext())

			if err != nil {
				return err
			}

			return ctx.JSON(users)
		}),
		view.Fragment(UserTableTarget, func(ctx fiber.Ctx) (templ.Component, error) {
			userService, err := NewUserService(ctx)

			if err != nil {
				return nil, err
			}

			users, err := userService.UserTableFragment(ctx.UserContext())

			if err != nil {
				return nil, err
			}

			return users, nil
		}),
		view.Page(func(ctx fiber.Ctx) (templ.Component, string, error) {
			userService, err := NewUserService(ctx)

			if err != nil {
				return nil, "", err
			}

			users, err := userService.UserTableFragment(ctx.UserContext())

			if err != nil {
				return nil, "", err
			}

			return users, "Manage system users", nil
		}),
	)

	app.Get("/new", view.Page(func(c fiber.Ctx) (page templ.Component, title string, err error) {
		return UserForm("/users", UserDTO{}), "Create new user", nil
	}))

	return app
}
