package utils

import (
	"fmt"

	"github.com/google/uuid"
)

func UserNotesCacheKey(userId uuid.UUID) string {
	return fmt.Sprintf("user:notes:v1:%s", userId.String())
}
func NoteByIdCacheKey(userId uuid.UUID, noteId uuid.UUID) string {
	return fmt.Sprintf("user:note:v1:%s:%s", userId.String(), noteId.String())
}
