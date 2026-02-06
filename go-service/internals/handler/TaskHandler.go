package handler

import (
	"encoding/json"
	middlewareV1 "go-servie/internals/middleware"
	"go-servie/utils"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

/* =========================
   CREATE TASK
========================= */

func (h *Handler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	user, err := middlewareV1.ExtractUser(r)
	if err != nil {
		utils.WriteJsonError(w, "Unauthorized", http.StatusUnauthorized, err)
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
		DueAt       string `json:"due_at"` // ISO string
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJsonError(w, "Invalid request body", http.StatusBadRequest, err)
		return
	}

	var dueAt time.Time
	if req.DueAt != "" {
		dueAt, err = time.Parse(time.RFC3339, req.DueAt)
		if err != nil {
			utils.WriteJsonError(w, "Invalid due_at format", http.StatusBadRequest, err)
			return
		}
	}

	task, err := h.service.CreateTask(
		r.Context(),
		user.UserId,
		req.Title,
		req.Description,
		req.Priority,
		dueAt,
	)
	if err != nil {
		utils.WriteJsonError(w, "Failed to create task", http.StatusInternalServerError, err)
		return
	}

	utils.WriteJsonESuccess(w, "Task created", http.StatusCreated, utils.TaskRes{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description.String,
		Status:      task.Status,
		Priority:    task.Priority,
		DueAt:       task.DueAt.Time,
		CompletedAt: task.CompletedAt.Time,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	})
}

/* =========================
   GET USER TASKS
========================= */

func (h *Handler) GetUserTasksHandler(w http.ResponseWriter, r *http.Request) {
	user, err := middlewareV1.ExtractUser(r)
	if err != nil {
		utils.WriteJsonError(w, "Unauthorized", http.StatusUnauthorized, err)
		return
	}

	tasks, err := h.service.GetUserTasks(r.Context(), user.UserId)
	if err != nil {
		utils.WriteJsonError(w, "Failed to fetch tasks", http.StatusInternalServerError, err)
		return
	}
	var result []utils.TaskRes
	for _, v := range tasks {
		result = append(result, utils.TaskRes{
			ID:          v.ID,
			UserID:      v.UserID,
			Title:       v.Title,
			Description: v.Description.String,
			Status:      v.Status,
			Priority:    v.Priority,
			DueAt:       v.DueAt.Time,
			CompletedAt: v.CompletedAt.Time,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}
	utils.WriteJsonESuccess(w, "Tasks fetched", http.StatusOK, result)
}

/* =========================
   GET TASK BY ID
========================= */

func (h *Handler) GetTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	user, err := middlewareV1.ExtractUser(r)
	if err != nil {
		utils.WriteJsonError(w, "Unauthorized", http.StatusUnauthorized, err)
		return
	}

	taskIDStr := chi.URLParam(r, "taskId")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		utils.WriteJsonError(w, "Invalid task id", http.StatusBadRequest, err)
		return
	}

	task, err := h.service.GetTaskById(
		r.Context(),
		taskID,
		user.UserId,
	)
	if err != nil {
		utils.WriteJsonError(w, "Task not found", http.StatusNotFound, err)
		return
	}

	utils.WriteJsonESuccess(w, "Task fetched", http.StatusOK, utils.TaskRes{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description.String,
		Status:      task.Status,
		Priority:    task.Priority,
		DueAt:       task.DueAt.Time,
		CompletedAt: task.CompletedAt.Time,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	})
}

/* =========================
   UPDATE TASK
========================= */

func (h *Handler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	user, err := middlewareV1.ExtractUser(r)
	if err != nil {
		utils.WriteJsonError(w, "Unauthorized", http.StatusUnauthorized, err)
		return
	}

	taskIDStr := chi.URLParam(r, "taskId")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		utils.WriteJsonError(w, "Invalid task id", http.StatusBadRequest, err)
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
		Status      string `json:"status"`
		DueAt       string `json:"due_at"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJsonError(w, "Invalid request body", http.StatusBadRequest, err)
		return
	}

	var dueAt time.Time
	if req.DueAt != "" {
		dueAt, err = time.Parse(time.RFC3339, req.DueAt)
		if err != nil {
			utils.WriteJsonError(w, "Invalid due_at format", http.StatusBadRequest, err)
			return
		}
	}

	task, err := h.service.UpdateTask(
		r.Context(),
		taskID,
		user.UserId,
		req.Title,
		req.Description,
		req.Priority,
		req.Status,
		dueAt,
	)
	if err != nil {
		utils.WriteJsonError(w, "Failed to update task", http.StatusInternalServerError, err)
		return
	}

	utils.WriteJsonESuccess(w, "Task updated", http.StatusOK, utils.TaskRes{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description.String,
		Status:      task.Status,
		Priority:    task.Priority,
		DueAt:       task.DueAt.Time,
		CompletedAt: task.CompletedAt.Time,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	})
}

/* =========================
   COMPLETE TASK
========================= */

func (h *Handler) CompleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	user, err := middlewareV1.ExtractUser(r)
	if err != nil {
		utils.WriteJsonError(w, "Unauthorized", http.StatusUnauthorized, err)
		return
	}

	taskIDStr := chi.URLParam(r, "taskId")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		utils.WriteJsonError(w, "Invalid task id", http.StatusBadRequest, err)
		return
	}

	task, err := h.service.CompleteTask(
		r.Context(),
		taskID,
		user.UserId,
	)
	if err != nil {
		utils.WriteJsonError(w, "Failed to complete task", http.StatusInternalServerError, err)
		return
	}

	utils.WriteJsonESuccess(w, "Task completed", http.StatusOK, utils.TaskRes{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description.String,
		Status:      task.Status,
		Priority:    task.Priority,
		DueAt:       task.DueAt.Time,
		CompletedAt: task.CompletedAt.Time,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	})
}

/* =========================
   DELETE TASK
========================= */

func (h *Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	user, err := middlewareV1.ExtractUser(r)
	if err != nil {
		utils.WriteJsonError(w, "Unauthorized", http.StatusUnauthorized, err)
		return
	}

	taskIDStr := chi.URLParam(r, "taskId")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		utils.WriteJsonError(w, "Invalid task id", http.StatusBadRequest, err)
		return
	}

	if err := h.service.DeleteTask(
		r.Context(),
		taskID,
		user.UserId,
	); err != nil {
		utils.WriteJsonError(w, "Failed to delete task", http.StatusInternalServerError, err)
		return
	}

	utils.WriteJsonESuccess(w, "Task deleted", http.StatusOK, map[string]any{
		"task_id": taskID,
	})
}
