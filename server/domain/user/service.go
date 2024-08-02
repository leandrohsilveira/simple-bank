package user_domain

import (
	"context"
	"fmt"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
	database "github.com/leandrohsilveira/simple-bank/server/database"
	"github.com/leandrohsilveira/simple-bank/server/domain"
	"github.com/leandrohsilveira/simple-bank/server/store"
)

const UserTableTarget = "users-table"
const UserFormTarget = "user-form"

type UserService struct {
	domain.DomainService
	store *store.Store
}

func NewUserService(ctx fiber.Ctx) (service UserService, err error) {
	service = UserService{}

	service.store, err = service.Store(ctx)

	return
}

func WithUserService(handler func(fiber.Ctx, UserService) error) fiber.Handler {
	return func(c fiber.Ctx) error {
		service, err := NewUserService(c)

		if err != nil {
			return err
		}

		return handler(c, service)
	}
}

func (service UserService) UserTableData(ctx context.Context) ([]UserDTO, error) {
	users, err := service.store.Queries.ListUsers(ctx, database.ListUsersParams{
		Offset: 0,
		Limit:  10,
	})

	if err != nil {
		return []UserDTO{}, err
	}

	return FromUsers(users), nil
}

func (service UserService) UserEditData(ctx context.Context, id int64) (UserDTO, error) {
	user, err := service.store.Queries.GetUserById(ctx, id)

	if err != nil {
		return UserDTO{}, err
	}

	return FromUser(user), nil
}

func (service UserService) UserTableFragment(ctx context.Context) (component templ.Component, err error) {
	users, err := service.UserTableData(ctx)

	if err != nil {
		return
	}

	component = UserTable(UserTableProps{
		Users:    users,
		Endpoint: "/users",
		Id:       UserTableTarget,
		Search:   "",
		Page:     1,
		Limit:    10,
		Count:    0,
	})

	return
}

func (service UserService) UserFormFragment(ctx context.Context, id *int64) (component templ.Component, err error) {
	if id == nil {
		return UserForm("/users/new", UserDTO{}), nil
	}

	user, err := service.UserEditData(ctx, *id)

	if err != nil {
		return nil, err
	}

	return UserForm(fmt.Sprintf("/users/%d", *id), user), nil
}
