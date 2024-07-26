package store

import (
	"context"
	"fmt"
	"log"
	"testing"

	database "github.com/leandrohsilveira/simple-bank/database/sqlc"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	
	ctx := context.Background()

	testStore, release, err := CreateStore(ctx)

	defer release()
	
	require.NoError(t, err)


	account1, _, err := database.CreateRandomAccount(ctx, testStore.Queries)
	
	require.NoError(t, err)
	
	account2, _, err := database.CreateRandomAccount(ctx, testStore.Queries)

	require.NoError(t, err)

	log.Printf("BEGIN BALANCES: from: %v, to: %v\n", account1.Balance, account2.Balance)

	errs := make(chan error)
	transfers := make(chan TransferResult)

	n := 10
	reverse := 3
	amountPerTransfer := 1000 / n

	require.LessOrEqual(t, reverse, n)
	require.IsType(t, int(0), amountPerTransfer)

	for i := 0; i < n; i++ {
		go func(index int) {
			fromAccount, toAccount := account1, account2

			if index < reverse {
				log.Printf("REVERSE TRANSFER %v\n", index)
				fromAccount, toAccount = account2, account1
			}

			arg := database.CreateTransferParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        int64(amountPerTransfer),
			}

			store, release, err := CreateStore(ctx)

			if err != nil {
				log.Printf("TRANSFER %v: error creating store: %v", index, err)
				release()
				transfers <- TransferResult{}
				errs <- fmt.Errorf("error creating store %v: %v", index, err)
				return
			}

			transfer, err := store.CreateTransfer(ctx, arg)

			if err != nil {
				log.Printf("TRANSFER %v: error creating transfer: %v", index, err)
				release()
				transfers <- TransferResult{}
				errs <- fmt.Errorf("error creating transfer %v: %v", index, err)
				return
			}

			transfers <- transfer

			release()

			log.Printf(
				"TRANSFER BALANCES %v: id: %v, acc1: %v, acc2: %v\n",
				index,
				transfer.Transfer.ID,
				transfer.FromAccount.Balance,
				transfer.ToAccount.Balance,
			)
		}(i)
	}

	for i := 0; i < n; i++ {
		transfer := <-transfers
		require.NotEmpty(t, transfer)
	
		require.Equal(t, database.TransferStatusCompleted, transfer.Transfer.Status)
	}

	close(errs)

	if len(errs) > 0 {
		errors := make([]error, 0)
		for err := range errs {
			errors = append(errors, err)
		}

		require.NoError(t, fmt.Errorf("errors: %v", errors))
	}

	account1Balance := account1.Balance
	account2Balance := account2.Balance

	log.Printf("FINAL BALANCES: from: %v, to: %v\n", account1.Balance, account2.Balance)

	account1, err = testStore.Queries.GetAccountById(ctx, account1.ID)
	require.NoError(t, err)

	account2, err = testStore.Queries.GetAccountById(ctx, account2.ID)
	require.NoError(t, err)

	reverseTransferedAmount := int64(amountPerTransfer * reverse)
	normalTransferedAmount := int64(amountPerTransfer * (n - reverse))

	require.Equal(t, account1Balance - normalTransferedAmount + reverseTransferedAmount, account1.Balance)
	require.Equal(t, account2Balance - reverseTransferedAmount + normalTransferedAmount, account2.Balance)
	
}