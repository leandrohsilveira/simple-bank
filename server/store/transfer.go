package store

import (
	"context"
	"fmt"
	"log"

	database "github.com/leandrohsilveira/simple-bank/server/database"
)

type TransferResult struct {
	Transfer    database.Transfer `json:"transfer"`
	FromAccount database.Account  `json:"from_account"`
	ToAccount   database.Account  `json:"to_account"`
	FromEntry   database.Entry    `json:"from_entry"`
	ToEntry     database.Entry    `json:"to_entry"`
}

func (store *Store) CreateTransfer(ctx context.Context, arg database.CreateTransferParams) (result TransferResult, err error) {
	result.Transfer, err = store.Queries.CreateTransfer(ctx, database.CreateTransferParams{
		FromAccountID: arg.FromAccountID,
		ToAccountID:   arg.ToAccountID,
		Amount:        arg.Amount,
		Status:        database.TransferStatusPending,
	})

	if err != nil {
		return TransferResult{}, err
	}

	err = store.ExecTx(ctx, func(q *database.Queries) error {
		return beginTransfer(ctx, q, &result)
	})

	if err != nil {
		err = failTransfer(ctx, store.Queries, result, err)

		return TransferResult{}, err
	}

	if result.Transfer.Status != database.TransferStatusTransfering {
		return result, nil
	}

	err = store.ExecTx(ctx, func(q *database.Queries) error {
		return completeTransfer(ctx, q, &result)
	})

	if err != nil {
		err = store.ExecTx(ctx, func(q *database.Queries) error {
			return failTransfer(ctx, store.Queries, result, err)
		})
		return TransferResult{}, err
	}

	return result, nil
}

func beginTransfer(ctx context.Context, q *database.Queries, result *TransferResult) (err error) {
	result.FromAccount, err = q.GetAccountForUpdate(ctx, result.Transfer.FromAccountID)

	if err != nil {
		return err
	}

	if result.FromAccount.Balance < result.Transfer.Amount {
		log.Printf(
			"insufficient from account balance %d for transfer ID %v amount %d",
			result.FromAccount.Balance,
			result.Transfer.ID,
			result.Transfer.Amount,
		)
		result.Transfer, err = rejectTransfer(ctx, q, result.Transfer.ID)

		return err
	}

	result.FromEntry, err = q.CreateEntry(ctx, database.CreateEntryParams{
		AccountID:  result.Transfer.FromAccountID,
		TransferID: result.Transfer.ID,
		Amount:     -result.Transfer.Amount,
	})

	if err != nil {
		return err
	}

	result.FromAccount, err = q.AddAccountBalance(ctx, database.AddAccountBalanceParams{
		ID:     result.Transfer.FromAccountID,
		Change: result.FromEntry.Amount,
	})

	if err != nil {
		return err
	}

	result.Transfer, err = q.SetTransferStatus(ctx, database.SetTransferStatusParams{
		ID:     result.Transfer.ID,
		Status: database.TransferStatusTransfering,
	})

	return err
}

func completeTransfer(ctx context.Context, q *database.Queries, result *TransferResult) (err error) {
	result.ToEntry, err = q.CreateEntry(ctx, database.CreateEntryParams{
		AccountID:  result.Transfer.ToAccountID,
		TransferID: result.Transfer.ID,
		Amount:     result.Transfer.Amount,
	})

	if err != nil {
		return err
	}

	result.ToAccount, err = q.AddAccountBalance(ctx, database.AddAccountBalanceParams{
		ID:     result.Transfer.ToAccountID,
		Change: result.ToEntry.Amount,
	})

	if err != nil {
		return err
	}

	result.Transfer, err = q.SetTransferStatus(ctx, database.SetTransferStatusParams{
		ID:     result.Transfer.ID,
		Status: database.TransferStatusCompleted,
	})

	return err
}

func failTransfer(ctx context.Context, q *database.Queries, result TransferResult, err error) error {
	var failedErr error

	_, failedErr = q.SetTransferStatus(ctx, database.SetTransferStatusParams{
		ID:     result.Transfer.ID,
		Status: database.TransferStatusFailed,
	})

	if failedErr != nil {
		return fmt.Errorf("failed to update transfer to failed: %v (caused by failed to transfer: %v)", failedErr, err)
	}

	if result.FromEntry != (database.Entry{}) {
		_, failedErr = q.AddAccountBalance(ctx, database.AddAccountBalanceParams{
			ID:     result.FromAccount.ID,
			Change: -result.FromEntry.Amount,
		})

		if failedErr != nil {
			return fmt.Errorf("failed to refund from account balance: %v (caused by failed to transfer: %v)", failedErr, err)
		}

		failedErr = q.DeleteEntry(ctx, result.FromEntry.ID)

		if failedErr != nil {
			return fmt.Errorf("failed to delete from entry: %v (caused by failed to transfer: %v)", failedErr, err)
		}
	}

	return err
}

func rejectTransfer(ctx context.Context, q *database.Queries, transferID int64) (database.Transfer, error) {
	return q.SetTransferStatus(ctx, database.SetTransferStatusParams{
		ID:     transferID,
		Status: database.TransferStatusRejected,
	})
}
