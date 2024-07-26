package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)


func TestCreateEntry(t *testing.T) {
	entry, arg, err := CreateRandomEntry(context.Background(), testQueries)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.CreatedAt)
}

func TestGetEntryById(t *testing.T) {
	createdEntry, _, err := CreateRandomEntry(context.Background(), testQueries)

	require.NoError(t, err)

	entry, err := testQueries.GetEntryById(context.Background(), createdEntry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, createdEntry.ID, entry.ID)
	require.Equal(t, createdEntry.AccountID, entry.AccountID)
	require.Equal(t, createdEntry.Amount, entry.Amount)
	require.Equal(t, createdEntry.CreatedAt, entry.CreatedAt)
}

func TestListEntries(t *testing.T) {
	createdEntry, _, err := CreateRandomEntry(context.Background(), testQueries)

	require.NoError(t, err)

	entries, err := testQueries.ListEntries(context.Background(), ListEntriesParams{
		Limit:  5,
		Offset: 0,
	})

	require.NoError(t, err)
	require.NotEmpty(t, entries)

	require.Equal(t, createdEntry.ID, entries[0].ID)
	require.Equal(t, createdEntry.AccountID, entries[0].AccountID)
	require.Equal(t, createdEntry.Amount, entries[0].Amount)
	require.Equal(t, createdEntry.CreatedAt, entries[0].CreatedAt)
}

func TestUpdateEntry(t *testing.T) {
	createdEntry, _, err := CreateRandomEntry(context.Background(), testQueries)

	require.NoError(t, err)

	arg := UpdateEntryParams{
		ID:     createdEntry.ID,
		Amount: 2000,
	}

	entry, err := testQueries.UpdateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, createdEntry.ID, entry.ID)
	require.Equal(t, createdEntry.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.Equal(t, createdEntry.CreatedAt, entry.CreatedAt)
}

func TestDeleteEntry(t *testing.T) {
	createdEntry, _, err := CreateRandomEntry(context.Background(), testQueries)

	require.NoError(t, err)

	err = testQueries.DeleteEntry(context.Background(), createdEntry.ID)

	require.NoError(t, err)

	_, err = testQueries.GetEntryById(context.Background(), createdEntry.ID)

	require.EqualError(t, err, "no rows in result set")
}

func TestCountEntries(t *testing.T) {
	_, _, err := CreateRandomEntry(context.Background(), testQueries)

	require.NoError(t, err)

	count, err := testQueries.CountEntries(context.Background())

	require.NoError(t, err)

	require.GreaterOrEqual(t, count, int64(1))
}