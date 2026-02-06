package repo

import (
	"context"
	"go-servie/dbmodel"

	"github.com/google/uuid"
)

func (searchrepo *Repo) SearchRepo(ctx context.Context, query string, userId uuid.UUID) ([]dbmodel.Note, error) {
	return searchrepo.db.SearchNotes(ctx, dbmodel.SearchNotesParams{
		PlaintoTsquery: query,
		UserID:         userId,
	})
}
