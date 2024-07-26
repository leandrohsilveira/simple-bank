package database

import (
	"context"

	"github.com/go-faker/faker/v4"
)

func CreateRandomAccount(ctx context.Context, queries *Queries) (Account, CreateAccountParams, error) {
	arg := CreateAccountParams{
		Owner:    faker.Name(),
		Balance:  1000,
		Currency: faker.Currency(),
	}

	account, err := queries.CreateAccount(ctx, arg)

	return account, arg, err
}

func CreateRandomTransfer(ctx context.Context, queries *Queries) (Transfer, CreateTransferParams, error) {
	fromAccount, _, err := CreateRandomAccount(ctx, queries)

	if err != nil {
		return Transfer{}, CreateTransferParams{}, err
	}

	toAccount, _, err := CreateRandomAccount(ctx, queries)

	if err != nil {
		return Transfer{}, CreateTransferParams{}, err
	}

	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        1000,
		Status:        TransferStatusPending,
	}

	transfer, err := queries.CreateTransfer(ctx, arg)

	return transfer, arg, err
}


func CreateRandomEntry(ctx context.Context, queries *Queries) (Entry, CreateEntryParams, error) {
	account, _, err := CreateRandomAccount(ctx, queries)

	if err != nil {
		return Entry{}, CreateEntryParams{}, err
	}

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    1000,
	}

	entry, err := queries.CreateEntry(ctx, arg)

	return entry, arg, err
}