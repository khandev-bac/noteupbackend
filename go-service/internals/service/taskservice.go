package service

import (
	"context"
	"errors"
	"go-servie/dbmodel"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidTaskInput = errors.New("invalid task input")
	ErrTaskNotFound     = errors.New("task not found")
)

/* =========================
   CREATE TASK
========================= */

func (s *Service) CreateTask(
	ctx context.Context,
	userID uuid.UUID,
	title string,
	description string,
	priority string,
	dueAt time.Time,
) (dbmodel.Task, error) {

	if userID == uuid.Nil || title == "" {
		return dbmodel.Task{}, ErrInvalidTaskInput
	}

	// default priority safety
	if priority == "" {
		priority = "medium"
	}

	return s.repo.CreateTaskRepo(
		ctx,
		userID,
		title,
		description,
		priority,
		dueAt,
	)
}

/* =========================
   GET ALL USER TASKS
========================= */

func (s *Service) GetUserTasks(
	ctx context.Context,
	userID uuid.UUID,
) ([]dbmodel.Task, error) {

	if userID == uuid.Nil {
		return nil, ErrInvalidTaskInput
	}

	return s.repo.GetUserTasksRepo(ctx, userID)
}

/* =========================
   GET TASK BY ID
========================= */

func (s *Service) GetTaskById(
	ctx context.Context,
	taskID uuid.UUID,
	userID uuid.UUID,
) (dbmodel.Task, error) {

	if taskID == uuid.Nil || userID == uuid.Nil {
		return dbmodel.Task{}, ErrInvalidTaskInput
	}

	return s.repo.GetTaskByIdRepo(ctx, taskID, userID)
}

/* =========================
   UPDATE TASK
========================= */

func (s *Service) UpdateTask(
	ctx context.Context,
	taskID uuid.UUID,
	userID uuid.UUID,
	title string,
	description string,
	priority string,
	status string,
	dueAt time.Time,
) (dbmodel.Task, error) {

	if taskID == uuid.Nil || userID == uuid.Nil {
		return dbmodel.Task{}, ErrInvalidTaskInput
	}

	if title == "" {
		return dbmodel.Task{}, ErrInvalidTaskInput
	}

	if status == "" {
		status = "pending"
	}

	if priority == "" {
		priority = "medium"
	}

	return s.repo.UpdateTaskRepo(
		ctx,
		taskID,
		userID,
		title,
		description,
		priority,
		status,
		dueAt,
	)
}

/* =========================
   COMPLETE TASK
========================= */

func (s *Service) CompleteTask(
	ctx context.Context,
	taskID uuid.UUID,
	userID uuid.UUID,
) (dbmodel.Task, error) {

	return s.repo.UpdateTaskRepo(
		ctx,
		taskID,
		userID,
		"",       // title unchanged (handled in repo/sql)
		"",       // description unchanged
		"medium", // priority unchanged
		"completed",
		time.Now(),
	)
}

/* =========================
   DELETE TASK
========================= */

func (s *Service) DeleteTask(
	ctx context.Context,
	taskID uuid.UUID,
	userID uuid.UUID,
) error {

	if taskID == uuid.Nil || userID == uuid.Nil {
		return ErrInvalidTaskInput
	}

	return s.repo.DeleteTaskRepo(ctx, taskID, userID)
}
