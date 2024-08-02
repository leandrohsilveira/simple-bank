package user_domain

import (
	"context"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
	database "github.com/leandrohsilveira/simple-bank/server/database"
	"github.com/leandrohsilveira/simple-bank/server/domain"
	"github.com/leandrohsilveira/simple-bank/server/store"
)

const UserTableTarget = "users-table"

type UserService struct {
	domain.DomainService
	store *store.Store
}

func NewUserService(ctx fiber.Ctx) (service UserService, err error) {
	service = UserService{}

	service.store, err = service.Store(ctx)

	return
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
