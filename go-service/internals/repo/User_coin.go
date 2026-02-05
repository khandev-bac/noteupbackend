package repo

import (
	"context"

	"go-servie/dbmodel"

	"github.com/google/uuid"
)

func (r *Repo) CreateUserCoinsRepo(ctx context.Context, userID uuid.UUID, balance int32) error {
	return r.db.CreateUserCoins(ctx, dbmodel.CreateUserCoinsParams{
		UserID:  userID,
		Balance: balance,
	})
}

func (r *Repo) GetUserCoinBalanceRepo(ctx context.Context, userID uuid.UUID) (int32, error) {
	return r.db.GetUserCoinBalance(ctx, userID)
}

func (r *Repo) GetUserCoinsRepo(ctx context.Context, userID uuid.UUID) (dbmodel.UserCoin, error) {
	return r.db.GetUserCoins(ctx, userID)
}

func (r *Repo) DeductUserCoinsRepo(ctx context.Context, userID uuid.UUID, coins int32) error {
	return r.db.DeductUserCoins(ctx, dbmodel.DeductUserCoinsParams{
		UserID:  userID,
		Balance: coins,
	})
}

func (r *Repo) AddUserCoinsRepo(ctx context.Context, userID uuid.UUID, coins int32) error {
	return r.db.AddUserCoins(ctx, dbmodel.AddUserCoinsParams{
		UserID:  userID,
		Balance: coins,
	})
}
