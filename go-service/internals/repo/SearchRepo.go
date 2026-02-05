package repo

import (
	"context"
	"go-servie/dbmodel"

	"github.com/google/uuid"
)

func (searchrepo *Repo) SearchRepo(ctx context.Context, userId uuid.UUID, query string) ([]dbmodel.Note, error) {
	return searchrepo.db.SearchNotes(ctx, dbmodel.SearchNotesParams{
		UserID:         userId,
		PlaintoTsquery: query,
	})
}
