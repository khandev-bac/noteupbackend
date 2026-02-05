package service

import (
	"context"
	"go-servie/dbmodel"

	"github.com/google/uuid"
)

func (searchservice *Service) SearchService(ctx context.Context, userId uuid.UUID, query string) ([]dbmodel.Note, error) {
	return searchservice.repo.SearchRepo(ctx, userId, query)
}
