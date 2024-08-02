package database

import (
	"context"
	"testing"

	"github.com/leandrohsilveira/simple-bank/server/database"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	transfer, arg, err := CreateRandomTransfer(context.Background(), testQueries)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.CreatedAt)
}

func TestGetTransferById(t *testing.T) {
	createdTransfer, _, err := CreateRandomTransfer(context.Background(), testQueries)

	require.NoError(t, err)

	transfer, err := testQueries.GetTransferById(context.Background(), createdTransfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, createdTransfer.ID, transfer.ID)
	require.Equal(t, createdTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, createdTransfer.ToAccountID, transfer.ToAccountID)
	require.Equal(t, createdTransfer.Amount, transfer.Amount)
	require.Equal(t, createdTransfer.CreatedAt, transfer.CreatedAt)
}

func TestListTransfers(t *testing.T) {
	createdTransfer, _, err := CreateRandomTransfer(context.Background(), testQueries)

	require.NoError(t, err)

	transfers, err := testQueries.ListTransfers(context.Background(), database.ListTransfersParams{
		Limit:  5,
		Offset: 0,
	})

	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	require.Equal(t, createdTransfer.ID, transfers[0].ID)
	require.Equal(t, createdTransfer.FromAccountID, transfers[0].FromAccountID)
	require.Equal(t, createdTransfer.ToAccountID, transfers[0].ToAccountID)
	require.Equal(t, createdTransfer.Amount, transfers[0].Amount)
	require.Equal(t, createdTransfer.CreatedAt, transfers[0].CreatedAt)
}

func TestUpdateTransfer(t *testing.T) {
	createdTransfer, _, err := CreateRandomTransfer(context.Background(), testQueries)

	require.NoError(t, err)

	arg := database.UpdateTransferParams{
		ID:     createdTransfer.ID,
		Amount: 1000,
		Status: database.TransferStatusPending,
	}

	transfer, err := testQueries.UpdateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, createdTransfer.ID, transfer.ID)
	require.Equal(t, createdTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, createdTransfer.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.Equal(t, createdTransfer.CreatedAt, transfer.CreatedAt)
}

func TestDeleteTransfer(t *testing.T) {
	createdTransfer, _, err := CreateRandomTransfer(context.Background(), testQueries)

	require.NoError(t, err)

	err = testQueries.DeleteTransfer(context.Background(), createdTransfer.ID)

	require.NoError(t, err)

	_, err = testQueries.GetTransferById(context.Background(), createdTransfer.ID)

	require.EqualError(t, err, "no rows in result set")
}

func TestCountTransfers(t *testing.T) {
	_, _, err := CreateRandomTransfer(context.Background(), testQueries)

	require.NoError(t, err)

	count, err := testQueries.CountTransfers(context.Background())

	require.NoError(t, err)
	require.GreaterOrEqual(t, count, int64(1))
}
