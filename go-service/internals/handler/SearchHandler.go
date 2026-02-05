package handler

import (
	"go-servie/utils"
	"net/http"
	"strings"
)

func (h *Handler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	// userId, err := middleware.ExtractUser(r)
	// if err != nil {
	// 	utils.WriteJsonError(w, "Unauthorized", http.StatusUnauthorized, err)
	// 	return
	// }
	ctx := r.Context()
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	query = strings.ToLower(query)
	if query == "" {
		utils.WriteJsonError(w, "Missing query", http.StatusBadRequest, nil)
		return
	}
	note, err := h.service.SearchService(ctx, query)
	if err != nil {
		utils.WriteJsonError(w, "Something went wrong in search", http.StatusBadRequest, nil)

	}
	var result []utils.NoteRes
	for _, v := range note {
		result = append(result, utils.NoteRes{
			ID:                   v.ID,
			UserID:               v.UserID,
			AudioUrl:             v.AudioUrl.String,
			AudioDurationSeconds: v.AudioDurationSeconds.Int32,
			AudioFileSizeMb:      v.AudioFileSizeMb.Int32,
			Title:                v.Title.String,
			Transcript:           v.Transcript.String,
			CreatedAt:            v.CreatedAt,
		})
	}
	utils.WriteJsonESuccess(w, "fetched", http.StatusOK, result)
}
