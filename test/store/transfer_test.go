package store

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/leandrohsilveira/simple-bank/server/database"
	"github.com/leandrohsilveira/simple-bank/server/store"
	test_database "github.com/leandrohsilveira/simple-bank/test/database"
	"github.com/stretchr/testify/require"
)

func TestCreateConcurrentTransfers(t *testing.T) {

	ctx := context.Background()

	testStore, release, err := CreateStore(ctx)

	require.NoError(t, err)

	defer release()

	account1, _, err := test_database.CreateRandomAccount(ctx, testStore.Queries)

	require.NoError(t, err)

	account2, _, err := test_database.CreateRandomAccount(ctx, testStore.Queries)

	require.NoError(t, err)

	log.Printf("BEGIN BALANCES: from: %v, to: %v\n", account1.Balance, account2.Balance)

	errs := make(chan error)
	transfers := make(chan store.TransferResult)

	n := 10
	reverse := 3
	amountPerTransfer := 1000 / n

	require.LessOrEqual(t, reverse, n)
	require.IsType(t, int(0), amountPerTransfer)

	for i := 0; i < n; i++ {
		go func(index int) {
			isReverse, fromAccount, toAccount, direction := false, account1, account2, "->"

			if index < reverse {
				isReverse, fromAccount, toAccount, direction = true, account2, account1, "<-"
			}

			arg := database.CreateTransferParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        int64(amountPerTransfer),
			}

			storeinstance, release, err := CreateStore(ctx)

			if err != nil {
				log.Printf("TRANSFER %v: error creating store: %v", index, err)
				transfers <- store.TransferResult{}
				errs <- fmt.Errorf("error creating store %v: %v", index, err)
				return
			}

			transfer, err := storeinstance.CreateTransfer(ctx, arg)

			release()

			if err != nil {
				log.Printf("TRANSFER %v: error creating transfer: %v", index, err)
				transfers <- store.TransferResult{}
				errs <- fmt.Errorf("error creating transfer %v: %v", index, err)
				return
			}

			transfers <- transfer

			acc1, acc2 := transfer.FromAccount, transfer.ToAccount

			if isReverse {
				acc1, acc2 = transfer.ToAccount, transfer.FromAccount
			}

			log.Printf(
				"TRANSFER BALANCES %v: id: %v, acc1: %v %v acc2: %v\n",
				index,
				transfer.Transfer.ID,
				acc1.Balance,
				direction,
				acc2.Balance,
			)
		}(i)
	}

	errors := make([]error, 0)

	for i := 0; i < n; i++ {
		transfer := <-transfers
		if transfer == (store.TransferResult{}) {
			errors = append(errors, <-errs)
		} else {
			require.Equal(t, database.TransferStatusCompleted, transfer.Transfer.Status)
		}
	}

	if len(errors) > 0 {
		require.Fail(t, "some transfers failed", errors)
	}

	account1Balance := account1.Balance
	account2Balance := account2.Balance

	account1, err = testStore.Queries.GetAccountById(ctx, account1.ID)
	require.NoError(t, err)

	account2, err = testStore.Queries.GetAccountById(ctx, account2.ID)
	require.NoError(t, err)

	log.Printf("FINAL BALANCES: from: %v, to: %v\n", account1.Balance, account2.Balance)

	reverseTransferedAmount := int64(amountPerTransfer * reverse)
	normalTransferedAmount := int64(amountPerTransfer * (n - reverse))

	require.Equal(t, account1Balance-normalTransferedAmount+reverseTransferedAmount, account1.Balance)
	require.Equal(t, account2Balance-reverseTransferedAmount+normalTransferedAmount, account2.Balance)

}

func TestCreateTransfersInsuficientFunds(t *testing.T) {

	ctx := context.Background()

	storeinstance, release, err := CreateStore(ctx)

	require.NoError(t, err)

	defer release()

	fromAccount, _, err := test_database.CreateRandomAccount(ctx, storeinstance.Queries)

	require.NoError(t, err)

	toAccount, _, err := test_database.CreateRandomAccount(ctx, storeinstance.Queries)

	require.NoError(t, err)

	n := 10
	totalAmount := 2000
	amountPerTransfer := totalAmount / n
	rejectedTransfersCount := (totalAmount - int(fromAccount.Balance)) / amountPerTransfer

	errs := make(chan error)
	transfers := make(chan store.TransferResult)

	require.IsType(t, int(0), amountPerTransfer)
	require.IsType(t, int(0), rejectedTransfersCount)
	require.Less(t, rejectedTransfersCount, n)

	for i := 0; i < n; i++ {
		go func(index int) {
			arg := database.CreateTransferParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        int64(amountPerTransfer),
			}

			storeinstance, release, err := CreateStore(ctx)

			if err != nil {
				log.Printf("TRANSFER %v: error creating store: %v", index, err)
				errs <- fmt.Errorf("error creating store %v: %v", index, err)
				transfers <- store.TransferResult{}
				return
			}

			transfer, err := storeinstance.CreateTransfer(ctx, arg)

			release()

			if err != nil {

				log.Printf("TRANSFER %v: error creating transfer: %v", index, err)
				errs <- fmt.Errorf("error creating transfer %v: %v", index, err)
				transfers <- store.TransferResult{}
				return
			}

			transfers <- transfer
		}(i)
	}

	errors := make([]error, 0)

	successfulTransfers := make([]store.TransferResult, 0)
	rejectedTransfers := make([]store.TransferResult, 0)

	for i := 0; i < n; i++ {
		transfer := <-transfers
		if transfer == (store.TransferResult{}) {
			errors = append(errors, <-errs)
		} else {
			require.Contains(t, []database.TransferStatus{database.TransferStatusCompleted, database.TransferStatusRejected}, transfer.Transfer.Status)

			switch transfer.Transfer.Status {
			case database.TransferStatusCompleted:
				successfulTransfers = append(successfulTransfers, transfer)
			case database.TransferStatusRejected:
				rejectedTransfers = append(rejectedTransfers, transfer)
			}
		}

	}

	if len(errors) > 0 {
		require.Fail(t, "some transfers failed", errors)
	}

	require.Equal(t, n-rejectedTransfersCount, len(successfulTransfers))
	require.Equal(t, rejectedTransfersCount, len(rejectedTransfers))
}
