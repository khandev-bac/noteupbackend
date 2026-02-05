package service

import (
	"context"

	"go-servie/dbmodel"

	"github.com/google/uuid"
)

func (s *Service) CreateCoinTransactionService(
	ctx context.Context,
	userID uuid.UUID,
	amount int32,
	reason string,
) error {
	return s.repo.CreateCoinTransactionRepo(ctx, userID, amount, reason)
}

func (s *Service) GetUserCoinTransactionsService(
	ctx context.Context,
	userID uuid.UUID,
) ([]dbmodel.CoinTransaction, error) {
	return s.repo.GetUserCoinTransactionsRepo(ctx, userID)
}
