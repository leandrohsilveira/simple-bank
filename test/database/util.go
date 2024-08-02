package database

import (
	"context"

	"github.com/go-faker/faker/v4"
	"github.com/leandrohsilveira/simple-bank/server/database"
)

func CreateRandomUser(ctx context.Context, queries *database.Queries) (user database.User, arg database.CreateUserParams, err error) {
	arg = database.CreateUserParams{
		Name:  faker.Name(),
		Email: faker.Email(),
		Role:  database.UserRoleRegularUser,
	}

	user, err = queries.CreateUser(ctx, arg)

	return
}

func CreateRandomAccount(ctx context.Context, queries *database.Queries) (database.Account, database.CreateAccountParams, error) {
	owner, _, err := CreateRandomUser(ctx, queries)

	if err != nil {
		return database.Account{}, database.CreateAccountParams{}, err
	}

	arg := database.CreateAccountParams{
		OwnerID:  owner.ID,
		Balance:  1000,
		Currency: faker.Currency(),
	}

	account, err := queries.CreateAccount(ctx, arg)

	return account, arg, err
}

func CreateRandomTransfer(ctx context.Context, queries *database.Queries) (database.Transfer, database.CreateTransferParams, error) {
	fromAccount, _, err := CreateRandomAccount(ctx, queries)

	if err != nil {
		return database.Transfer{}, database.CreateTransferParams{}, err
	}

	toAccount, _, err := CreateRandomAccount(ctx, queries)

	if err != nil {
		return database.Transfer{}, database.CreateTransferParams{}, err
	}

	arg := database.CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        1000,
		Status:        database.TransferStatusPending,
	}

	transfer, err := queries.CreateTransfer(ctx, arg)

	return transfer, arg, err
}

func CreateRandomEntry(ctx context.Context, queries *database.Queries) (database.Entry, database.CreateEntryParams, error) {
	transfer, _, err := CreateRandomTransfer(ctx, queries)

	if err != nil {
		return database.Entry{}, database.CreateEntryParams{}, err
	}

	arg := database.CreateEntryParams{
		AccountID:  transfer.ToAccountID,
		TransferID: transfer.ID,
		Amount:     1000,
	}

	entry, err := queries.CreateEntry(ctx, arg)

	return entry, arg, err
}
