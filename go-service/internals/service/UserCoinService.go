package service

import (
	"context"

	"go-servie/dbmodel"

	"github.com/google/uuid"
)

/* =========================
   USER COINS SERVICE
========================= */

func (s *Service) CreateUserCoinsService(ctx context.Context, userID uuid.UUID, balance int32) error {
	return s.repo.CreateUserCoinsRepo(ctx, userID, balance)
}

func (s *Service) GetUserCoinBalanceService(ctx context.Context, userID uuid.UUID) (int32, error) {
	return s.repo.GetUserCoinBalanceRepo(ctx, userID)
}

func (s *Service) GetUserCoinsService(ctx context.Context, userID uuid.UUID) (*dbmodel.UserCoin, error) {
	coins, err := s.repo.GetUserCoinsRepo(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &coins, nil
}

func (s *Service) DeductUserCoinsService(ctx context.Context, userID uuid.UUID, coins int32) error {
	return s.repo.DeductUserCoinsRepo(ctx, userID, coins)
}

func (s *Service) AddUserCoinsService(ctx context.Context, userID uuid.UUID, coins int32) error {
	return s.repo.AddUserCoinsRepo(ctx, userID, coins)
}
