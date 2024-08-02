package user_domain

import (
	"github.com/gofiber/fiber/v3"
	"github.com/leandrohsilveira/simple-bank/server/view"
)

func UserRouter() *fiber.App {
	app := fiber.New()

	app.Get("/",
		view.Json(WithUserService(func(ctx fiber.Ctx, userService UserService) error {
			users, err := userService.UserTableData(ctx.UserContext())

			if err != nil {
				return err
			}

			return ctx.JSON(users)
		})),
		view.Fragment(UserTableTarget, WithUserService(func(ctx fiber.Ctx, userService UserService) error {
			table, err := userService.UserTableFragment(ctx.UserContext())

			if err != nil {
				return err
			}

			return view.Render(ctx, table)
		})),
		view.Page(WithUserService(func(ctx fiber.Ctx, userService UserService) error {
			table, err := userService.UserTableFragment(ctx.UserContext())

			if err != nil {
				return err
			}

			return view.RenderPage(ctx, table, "Manage system users")
		})),
	)

	app.Get("/new", view.Page(func(ctx fiber.Ctx) error {
		return view.RenderPage(ctx, UserForm("/users", UserDTO{}), "Create new user")
	}))

	app.Get(
		"/:id",
		view.Json(WithUserService(func(ctx fiber.Ctx, userService UserService) error {
			id := fiber.Params[int64](ctx, "id")

			user, err := userService.UserEditData(ctx.UserContext(), id)

			if err != nil {
				return err
			}

			return ctx.JSON(user)
		})),
		view.Fragment(UserFormTarget, WithUserService(func(ctx fiber.Ctx, userService UserService) error {
			id := fiber.Params[int64](ctx, "id")

			form, err := userService.UserFormFragment(ctx.UserContext(), &id)

			if err != nil {
				return err
			}

			return view.Render(ctx, form)

		})),
		view.Page(WithUserService(func(ctx fiber.Ctx, userService UserService) error {
			id := fiber.Params[int64](ctx, "id")

			form, err := userService.UserFormFragment(ctx.UserContext(), &id)

			if err != nil {
				return err
			}

			return view.RenderPage(ctx, form, "Edit user")
		})),
	)

	return app
}
