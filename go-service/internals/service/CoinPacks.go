package service

import (
	"context"

	"go-servie/dbmodel"

	"github.com/google/uuid"
)

func (s *Service) GetActiveCoinPacksService(ctx context.Context) ([]dbmodel.GetActiveCoinPacksRow, error) {
	return s.repo.GetActiveCoinPacksRepo(ctx)
}

func (s *Service) GetCoinPackByIdService(ctx context.Context, packID uuid.UUID) (*dbmodel.CoinPack, error) {
	pack, err := s.repo.GetCoinPackByIdRepo(ctx, packID)
	if err != nil {
		return nil, err
	}
	return &pack, nil
}
