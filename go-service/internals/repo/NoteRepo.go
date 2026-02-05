package repo

import (
	"context"
	"database/sql"
	"go-servie/dbmodel"

	"github.com/google/uuid"
)

func (noterepo *Repo) CreateNoteRepo(ctx context.Context, user_id uuid.UUID, audio_url string, filesize int32, audio_duration int32) (dbmodel.Note, error) {
	return noterepo.db.CreateNotes(ctx, dbmodel.CreateNotesParams{
		UserID: user_id,
		AudioUrl: sql.NullString{
			String: audio_url,
			Valid:  true,
		},
		AudioFileSizeMb:      sql.NullInt32{Int32: filesize, Valid: true},
		AudioDurationSeconds: sql.NullInt32{Int32: audio_duration, Valid: true},
	})
}

func (noterepo *Repo) AfterProcessingUpdateNotesRepo(ctx context.Context, note_id uuid.UUID, title, transcript string, word_count int32) (dbmodel.Note, error) {
	return noterepo.db.AfterProcessingUpdateNotes(ctx, dbmodel.AfterProcessingUpdateNotesParams{
		ID:         note_id,
		Title:      sql.NullString{String: title, Valid: true},
		Transcript: sql.NullString{String: transcript, Valid: true},
		WordCount:  sql.NullInt32{Int32: word_count, Valid: true},
	})
}
func (noterepo *Repo) UpdateNoteWithNoteIdRepo(ctx context.Context, note_id uuid.UUID, title, transcript string) (dbmodel.Note, error) {
	return noterepo.db.UpdateNoteWithNoteId(ctx, dbmodel.UpdateNoteWithNoteIdParams{
		ID:         note_id,
		Title:      sql.NullString{String: title, Valid: true},
		Transcript: sql.NullString{String: transcript, Valid: true},
	})
}
func (r *Repo) GetAllUsersNotesRepo(
	ctx context.Context,
	userID uuid.UUID,
) ([]dbmodel.Note, error) {
	return r.db.GetAllUsersNotes(ctx, userID)
}

func (r *Repo) GetNoteByNoteIdRepo(
	ctx context.Context,
	noteID uuid.UUID,
	userID uuid.UUID,
) (dbmodel.Note, error) {
	return r.db.GetNoteByNoteId(ctx, dbmodel.GetNoteByNoteIdParams{
		ID:     noteID,
		UserID: userID,
	})
}
func (r *Repo) DeleteNoteByIdRepo(
	ctx context.Context,
	noteID uuid.UUID,
	userID uuid.UUID,
) error {
	return r.db.DeleteNoteById(ctx, dbmodel.DeleteNoteByIdParams{
		ID:     noteID,
		UserID: userID,
	})
}
