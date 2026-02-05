package service

import (
	"context"
	"go-servie/dbmodel"
)

func (searchservice *Service) SearchService(ctx context.Context, query string) ([]dbmodel.Note, error) {
	return searchservice.repo.SearchRepo(ctx, query)
}
