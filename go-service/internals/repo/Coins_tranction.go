package repo

import (
	"context"

	"go-servie/dbmodel"

	"github.com/google/uuid"
)

func (r *Repo) CreateCoinTransactionRepo(
	ctx context.Context,
	userID uuid.UUID,
	amount int32,
	reason string,
) error {
	return r.db.CreateCoinTransaction(ctx, dbmodel.CreateCoinTransactionParams{
		UserID: userID,
		Amount: amount,
		Reason: reason,
	})
}

func (r *Repo) GetUserCoinTransactionsRepo(
	ctx context.Context,
	userID uuid.UUID,
) ([]dbmodel.CoinTransaction, error) {
	return r.db.GetUserCoinTransactions(ctx, userID)
}
