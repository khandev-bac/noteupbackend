package handler

import (
	"encoding/json"
	"fmt"
	"go-servie/internals/config"
	middlewareV1 "go-servie/internals/middleware"
	"go-servie/utils"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/imagekit-developer/imagekit-go/v2"
)

func (h *Handler) CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	user, err := middlewareV1.ExtractUser(r)
	if err != nil {
		utils.WriteJsonError(w, "Unauthorized", http.StatusUnauthorized, err)
		return
	}

	// duration from frontend
	durationStr := r.FormValue("duration_seconds")
	if durationStr == "" {
		utils.WriteJsonError(w, "missing duration_seconds", http.StatusBadRequest, nil)
		return
	}

	durationSeconds, err := strconv.Atoi(durationStr)
	if err != nil || durationSeconds <= 0 {
		utils.WriteJsonError(w, "invalid duration_seconds", http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	// coin check
	balance, err := h.service.GetUserCoinBalanceService(ctx, user.UserId)
	if err != nil {
		utils.WriteJsonError(w, "failed to fetch coin balance", 500, err)
		return
	}

	const SecondsPerCoin = 30
	maxAllowed := balance * SecondsPerCoin

	if int32(durationSeconds) > maxAllowed {
		utils.WriteJsonError(
			w,
			fmt.Sprintf("max allowed %d seconds", maxAllowed),
			http.StatusPaymentRequired,
			nil,
		)
		return
	}

	// parse audio
	r.Body = http.MaxBytesReader(w, r.Body, 25<<20)
	if err := r.ParseMultipartForm(25 << 20); err != nil {
		utils.WriteJsonError(w, "invalid form", 400, err)
		return
	}

	file, header, err := r.FormFile("audio")
	if err != nil {
		utils.WriteJsonError(w, "audio missing", 400, err)
		return
	}
	defer file.Close()

	// upload to imagekit
	ik := config.ImagekitConfig()
	resp, err := ik.Files.Upload(ctx, imagekit.FileUploadParams{
		File:     file,
		FileName: header.Filename,
	})
	if err != nil {
		utils.WriteJsonError(w, "upload failed", 500, err)
		return
	}

	// create note (processing)
	note, err := h.service.CreateNote(
		ctx,
		user.UserId,
		resp.URL,
		int32(resp.Size),
		int32(durationSeconds),
	)
	if err != nil {
		utils.WriteJsonError(w, "db insert failed", 500, err)
		return
	}

	// deduct coins NOW (safe)
	coins := int32((durationSeconds + SecondsPerCoin - 1) / SecondsPerCoin)
	_ = h.redis.Del(ctx, utils.UserNotesCacheKey(user.UserId)).Err()
	_ = h.service.DeductUserCoinsService(ctx, user.UserId, coins)
	_ = h.service.CreateCoinTransactionService(
		ctx,
		user.UserId,
		coins,
		"note_creation",
	)

	// publish kafka
	err = h.kafkaConfig.PublishNoteCreated(
		utils.NOTECREATEDEVENT,
		note.ID.String(),
		resp.URL,
	)
	if err != nil {
		utils.WriteJsonError(w, "kafka failed", 500, err)
		return
	}

	// 🔥 RETURN IMMEDIATELY (no waiting)
	utils.WriteJsonESuccess(w, "processing started", http.StatusAccepted, map[string]any{
		"note_id": note.ID,
		"status":  "processing",
	})
}

func (h *Handler) UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	user, err := middlewareV1.ExtractUser(r)
	if err != nil {
		utils.WriteJsonError(w, "Unauthorized", http.StatusUnauthorized, err)
		return
	}

	ctx := r.Context()
	var req utils.NoteUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJsonError(w, "Invalid request", http.StatusBadRequest, err)
		return
	}

	noteIdstr := chi.URLParam(r, "noteId")
	noteId, err := uuid.Parse(noteIdstr)
	fmt.Println("Error while parsing uuid: ", err)

	note, err := h.service.UpdateNote(ctx, noteId, req.Title, req.Transcript)
	if err != nil {
		utils.WriteJsonError(w, "Failed to update notes", http.StatusInternalServerError, err)
		return
	}

	// 🔥 CACHE INVALIDATION (ONLY ADDITION)
	_ = h.redis.Del(
		ctx,
		utils.UserNotesCacheKey(user.UserId),
		utils.NoteByIdCacheKey(user.UserId, noteId),
	).Err()

	utils.WriteJsonESuccess(w, "Successfully updated", http.StatusOK, utils.NoteRes{
		ID:                   note.ID,
		UserID:               note.UserID,
		AudioUrl:             note.AudioUrl.String,
		AudioDurationSeconds: note.AudioDurationSeconds.Int32,
		AudioFileSizeMb:      note.AudioDurationSeconds.Int32,
		Transcript:           note.Transcript.String,
		Title:                note.Title.String,
		WordCount:            note.WordCount.Int32,
		SearchVector:         note.SearchVector,
		Status:               note.Status,
		CreatedAt:            note.CreatedAt,
		UpdatedAt:            note.UpdatedAt,
	})
}

func (h *Handler) UsersNotesHandler(w http.ResponseWriter, r *http.Request) {
	User, err := middlewareV1.ExtractUser(r)
	if err != nil {
		utils.WriteJsonError(w, "Unauthorized", http.StatusUnauthorized, err)
		return
	}
	ctx := r.Context()
	cacheKey := utils.UserNotesCacheKey(User.UserId)

	// 🔹 CACHE HIT
	if cached, err := h.redis.Get(ctx, cacheKey).Result(); err == nil {
		var res []utils.NoteRes
		log.Println("from cache all notes")
		if json.Unmarshal([]byte(cached), &res) == nil {
			utils.WriteJsonESuccess(w, "Fetched (cache)", http.StatusOK, res)
			return
		}
	}

	note, err := h.service.GetUserNotes(ctx, User.UserId)
	if err != nil {
		utils.WriteJsonError(w, "Failed to fetch notes", http.StatusInternalServerError, err)
		return
	}
	var result []utils.NoteRes
	for _, v := range note {
		result = append(result, utils.NoteRes{
			ID:                   v.ID,
			UserID:               v.UserID,
			AudioUrl:             v.AudioUrl.String,
			AudioDurationSeconds: v.AudioDurationSeconds.Int32,
			AudioFileSizeMb:      v.AudioDurationSeconds.Int32,
			Transcript:           v.Transcript.String,
			Title:                v.Title.String,
			WordCount:            v.WordCount.Int32,
			SearchVector:         v.SearchVector,
			Status:               v.Status,
			CreatedAt:            v.CreatedAt,
			UpdatedAt:            v.UpdatedAt,
		})
	}
	if data, err := json.Marshal(result); err == nil {
		_ = h.redis.SetEx(ctx, cacheKey, data, 30*time.Second).Err()
	}

	utils.WriteJsonESuccess(w, "Successfully fetched all notes", http.StatusOK, result)
}
func (h *Handler) NotesById(w http.ResponseWriter, r *http.Request) {
	User, err := middlewareV1.ExtractUser(r)
	if err != nil {
		utils.WriteJsonError(w, "Unauthorized", http.StatusUnauthorized, err)
		return
	}

	ctx := r.Context()
	noteId, err := uuid.Parse(chi.URLParam(r, "noteId"))
	if err != nil {
		utils.WriteJsonError(w, "Invalid note id", http.StatusBadRequest, err)
		return
	}

	cacheKey := utils.NoteByIdCacheKey(User.UserId, noteId)

	if cached, err := h.redis.Get(ctx, cacheKey).Result(); err == nil {
		var res utils.NoteRes
		log.Println("from cache notebyid")
		if json.Unmarshal([]byte(cached), &res) == nil {
			utils.WriteJsonESuccess(w, "Fetched (cache)", http.StatusOK, res)
			return
		}
	}

	note, err := h.service.GetNote(ctx, noteId, User.UserId)
	if err != nil {
		utils.WriteJsonError(w, "Failed to fetch note", http.StatusInternalServerError, err)
		return
	}

	res := utils.NoteRes{
		ID:                   note.ID,
		UserID:               note.UserID,
		AudioUrl:             note.AudioUrl.String,
		AudioDurationSeconds: note.AudioDurationSeconds.Int32,
		AudioFileSizeMb:      note.AudioDurationSeconds.Int32,
		Transcript:           note.Transcript.String,
		Title:                note.Title.String,
		WordCount:            note.WordCount.Int32,
		SearchVector:         note.SearchVector,
		Status:               note.Status,
		CreatedAt:            note.CreatedAt,
		UpdatedAt:            note.UpdatedAt,
	}

	if data, err := json.Marshal(res); err == nil {
		_ = h.redis.SetEx(ctx, cacheKey, data, 30*time.Second).Err()
	}

	utils.WriteJsonESuccess(w, "Fetched", http.StatusOK, res)
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	user, err := middlewareV1.ExtractUser(r)
	if err != nil {
		utils.WriteJsonError(w, "Unauthorized", http.StatusUnauthorized, err)
		return
	}

	ctx := r.Context()
	noteId, err := uuid.Parse(chi.URLParam(r, "noteId"))
	if err != nil {
		utils.WriteJsonError(w, "Invalid note id", http.StatusBadRequest, err)
		return
	}

	if err := h.service.DeleteNote(ctx, noteId, user.UserId); err != nil {
		utils.WriteJsonError(w, "Failed to delete", http.StatusInternalServerError, err)
		return
	}

	// 🔥 INVALIDATE
	_ = h.redis.Del(
		ctx,
		utils.UserNotesCacheKey(user.UserId),
		utils.NoteByIdCacheKey(user.UserId, noteId),
	).Err()

	utils.WriteJsonESuccess(w, "Deleted", http.StatusOK, map[string]any{
		"noteId": noteId,
	})
}

func (h *Handler) TestNote(w http.ResponseWriter, r *http.Request) {

	noteIdstr := chi.URLParam(r, "noteId")
	noteId, err := uuid.Parse(noteIdstr)
	fmt.Println("Error parse uuid: ", err)
	utils.WriteJsonESuccess(w, "Successfully printed", http.StatusOK, map[string]any{
		"noteId": noteId,
	})
}

func (h *Handler) TestAudioDuration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ik := config.ImagekitConfig()
	file, header, err := r.FormFile("audio")
	if err != nil {
		utils.WriteJsonError(w, "very short audio", http.StatusBadRequest, err)
		return
	}
	defer file.Close()
	resp, err := ik.Files.Upload(ctx, imagekit.FileUploadParams{
		File:     file,
		FileName: header.Filename,
	})
	if err != nil {
		utils.WriteJsonError(w, "failed to processes audio", http.StatusInternalServerError, err)
		return
	}
	utils.WriteJsonESuccess(
		w,
		"Audio test",
		http.StatusPaymentRequired,
		map[string]any{
			"audio_size":     resp.Size,
			"audio_duration": resp.Duration,
			"audio_type":     resp.FileType,
		},
	)
}
