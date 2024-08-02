package database

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leandrohsilveira/simple-bank/server/database"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	user, arg, err := CreateRandomUser(context.Background(), testQueries)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Role, user.Role)
	require.NotZero(t, user.CreatedAt)
}

func TestUpdateUser(t *testing.T) {
	createdUser, _, err := CreateRandomUser(context.Background(), testQueries)

	require.NoError(t, err)
	require.NotEmpty(t, createdUser)

	arg := database.UpdateUserParams{
		ID:    createdUser.ID,
		Email: faker.Email(),
		Name:  faker.Name(),
		Role:  database.UserRoleAdminUser,
	}

	user, err := testQueries.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, createdUser.ID, user.ID)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.Role, user.Role)
	require.Equal(t, createdUser.CreatedAt, user.CreatedAt)
}

func TestUpdateUserPassword(t *testing.T) {
	createdUser, _, err := CreateRandomUser(context.Background(), testQueries)

	require.NoError(t, err)
	require.NotEmpty(t, createdUser)

	arg := database.UpdateUserPasswordParams{
		ID:       createdUser.ID,
		Password: pgtype.Text{String: faker.Password(), Valid: true},
	}

	user, err := testQueries.UpdateUserPassword(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, createdUser.ID, user.ID)
	require.NotEmpty(t, arg.Password.String)
	require.NotEmpty(t, user.Password.String)
	require.NotEqual(t, createdUser.Password.String, user.Password.String)
	require.Equal(t, arg.Password.String, user.Password.String)
	require.Equal(t, createdUser.CreatedAt, user.CreatedAt)
}

func TestGetUserById(t *testing.T) {
	createdUser, _, err := CreateRandomUser(context.Background(), testQueries)

	require.NoError(t, err)

	user, err := testQueries.GetUserById(context.Background(), createdUser.ID)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, createdUser.ID, user.ID)
	require.Equal(t, createdUser.Name, user.Name)
	require.Equal(t, createdUser.Email, user.Email)
	require.Equal(t, createdUser.Role, user.Role)
}

func TestListUsers(t *testing.T) {
	createdUser, _, err := CreateRandomUser(context.Background(), testQueries)

	require.NoError(t, err)

	users, err := testQueries.ListUsers(context.Background(), database.ListUsersParams{
		Limit:  5,
		Offset: 0,
	})

	require.NoError(t, err)
	require.NotEmpty(t, users)

	require.Equal(t, createdUser.ID, users[0].ID)
	require.Equal(t, createdUser.Name, users[0].Name)
	require.Equal(t, createdUser.Email, users[0].Email)
	require.Equal(t, createdUser.Role, users[0].Role)
}

func TestDeleteUser(t *testing.T) {
	createdUser, _, err := CreateRandomUser(context.Background(), testQueries)

	require.NoError(t, err)

	err = testQueries.DeleteUser(context.Background(), createdUser.ID)

	require.NoError(t, err)

	_, err = testQueries.GetUserById(context.Background(), createdUser.ID)

	require.EqualError(t, err, "no rows in result set")
}
