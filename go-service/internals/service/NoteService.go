package service

import (
	"context"
	"errors"
	"go-servie/dbmodel"

	"github.com/google/uuid"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("note not found")
)

func (s *Service) CreateNote(
	ctx context.Context,
	userID uuid.UUID,
	audioURL string,
	filesize int32,
	audio_duration int32,
) (dbmodel.Note, error) {
	if userID == uuid.Nil {
		return dbmodel.Note{}, ErrInvalidInput
	}

	return s.repo.CreateNoteRepo(ctx, userID, audioURL, filesize, audio_duration)
}

func (s *Service) CompleteProcessing(
	ctx context.Context,
	noteID uuid.UUID,
	title string,
	transcript string,
	wordCount int32,
) (dbmodel.Note, error) {
	if noteID == uuid.Nil {
		return dbmodel.Note{}, ErrInvalidInput
	}

	return s.repo.AfterProcessingUpdateNotesRepo(
		ctx,
		noteID,
		title,
		transcript,
		wordCount,
	)
}

func (s *Service) UpdateNote(
	ctx context.Context,
	noteID uuid.UUID,
	title string,
	transcript string,
) (dbmodel.Note, error) {
	if noteID == uuid.Nil {
		return dbmodel.Note{}, ErrInvalidInput
	}

	return s.repo.UpdateNoteWithNoteIdRepo(
		ctx,
		noteID,
		title,
		transcript,
	)
}

func (s *Service) GetUserNotes(
	ctx context.Context,
	userID uuid.UUID,
) ([]dbmodel.Note, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidInput
	}

	return s.repo.GetAllUsersNotesRepo(ctx, userID)
}

func (s *Service) GetNote(
	ctx context.Context,
	noteID uuid.UUID,
	userID uuid.UUID,
) (dbmodel.Note, error) {
	if noteID == uuid.Nil || userID == uuid.Nil {
		return dbmodel.Note{}, ErrInvalidInput
	}

	return s.repo.GetNoteByNoteIdRepo(ctx, noteID, userID)
}

func (s *Service) DeleteNote(
	ctx context.Context,
	noteID uuid.UUID,
	userID uuid.UUID,
) error {
	if noteID == uuid.Nil || userID == uuid.Nil {
		return ErrInvalidInput
	}

	return s.repo.DeleteNoteByIdRepo(ctx, noteID, userID)
}
