package repo

import (
	"context"
	"go-servie/dbmodel"
)

func (searchrepo *Repo) SearchRepo(ctx context.Context, query string) ([]dbmodel.Note, error) {
	return searchrepo.db.SearchNotes(ctx, query)
}
