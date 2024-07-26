package store

import (
	"context"
	"fmt"

	database "github.com/leandrohsilveira/simple-bank/database/sqlc"
)

type TransferResult struct {
	Transfer    database.Transfer `json:"transfer"`
	FromAccount database.Account  `json:"from_account"`
	ToAccount   database.Account  `json:"to_account"`
	FromEntry   database.Entry    `json:"from_entry"`
	ToEntry     database.Entry    `json:"to_entry"`
}

func (store *Store) CreateTransfer(ctx context.Context, arg database.CreateTransferParams) (TransferResult, error) {
	var result TransferResult

	var err error

	err = store.ExecTx(ctx, func(q *database.Queries) error {
		
		result.Transfer, err = store.Queries.CreateTransfer(ctx, database.CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
			Status:        database.TransferStatusPending,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, database.CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, database.CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		return nil
	})

	if (err != nil) {
		return TransferResult{}, err
	}


	err = store.ExecTx(ctx, func(q *database.Queries) error {
		result.FromAccount, err = q.AddAccountBalance(ctx, database.AddAccountBalanceParams{
			ID:      arg.FromAccountID,
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
	})

	if (err != nil) {
		err = setTransferFailed(ctx, store.Queries, result.Transfer.ID, err)

		return TransferResult{}, err
	}

	err = store.ExecTx(ctx, func(q *database.Queries) error {
		result.ToAccount, err = q.AddAccountBalance(ctx, database.AddAccountBalanceParams{
			ID:      arg.ToAccountID,
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
	})

	if (err != nil) {
		err = setTransferFailed(ctx, store.Queries, result.Transfer.ID, err)
		return TransferResult{}, err
	}

	return result, nil
}

func setTransferFailed(ctx context.Context, q *database.Queries, transferID int64, err error) error {
	var failedErr error
	_, failedErr = q.SetTransferStatus(ctx, database.SetTransferStatusParams{
		ID:     transferID,
		Status: database.TransferStatusFailed,
	})

	if failedErr != nil {
		return fmt.Errorf("failed to update transfer to failed: %v (caused by failed to transfer: %v)", failedErr, err)
	}

	return err
}