package database

import (
	"context"
	"testing"

	"github.com/leandrohsilveira/simple-bank/server/database"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	account, arg, err := CreateRandomAccount(context.Background(), testQueries)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.OwnerID, account.OwnerID)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.CreatedAt)
}

func TestGetAccountById(t *testing.T) {
	createdAccount, _, err := CreateRandomAccount(context.Background(), testQueries)

	require.NoError(t, err)

	account, err := testQueries.GetAccountById(context.Background(), createdAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, createdAccount.ID, account.ID)
	require.Equal(t, createdAccount.OwnerID, account.OwnerID)
	require.Equal(t, createdAccount.Balance, account.Balance)
	require.Equal(t, createdAccount.Currency, account.Currency)
}

func TestListAccounts(t *testing.T) {
	createdAccount, _, err := CreateRandomAccount(context.Background(), testQueries)

	require.NoError(t, err)

	accounts, err := testQueries.ListAccounts(context.Background(), database.ListAccountsParams{
		Limit:  5,
		Offset: 0,
	})

	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	require.Equal(t, createdAccount.ID, accounts[0].ID)
	require.Equal(t, createdAccount.OwnerID, accounts[0].OwnerID)
	require.Equal(t, createdAccount.Balance, accounts[0].Balance)
	require.Equal(t, createdAccount.Currency, accounts[0].Currency)
}

func TestAddAccountBalance(t *testing.T) {
	createdAccount, _, err := CreateRandomAccount(context.Background(), testQueries)

	require.NoError(t, err)

	arg := database.AddAccountBalanceParams{
		ID:     createdAccount.ID,
		Change: 1000,
	}

	account, err := testQueries.AddAccountBalance(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, createdAccount.ID, account.ID)
	require.Equal(t, createdAccount.OwnerID, account.OwnerID)
	require.Equal(t, createdAccount.Balance+arg.Change, account.Balance)
	require.Equal(t, createdAccount.Currency, account.Currency)
}

func TestDeleteAccount(t *testing.T) {
	createdAccount, _, err := CreateRandomAccount(context.Background(), testQueries)

	require.NoError(t, err)

	err = testQueries.DeleteAccount(context.Background(), createdAccount.ID)

	require.NoError(t, err)

	_, err = testQueries.GetAccountById(context.Background(), createdAccount.ID)

	require.EqualError(t, err, "no rows in result set")
}

func TestCountAccounts(t *testing.T) {
	_, _, err := CreateRandomAccount(context.Background(), testQueries)

	require.NoError(t, err)

	count, err := testQueries.CountAccounts(context.Background())

	require.NoError(t, err)

	require.GreaterOrEqual(t, count, int64(1))
}
