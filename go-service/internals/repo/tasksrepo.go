package repo

import (
	"context"
	"database/sql"
	"go-servie/dbmodel"
	"time"

	"github.com/google/uuid"
)

func (r *Repo) CreateTaskRepo(
	ctx context.Context,
	userID uuid.UUID,
	title string,
	description string,
	priority string,
	dueAt time.Time,
) (dbmodel.Task, error) {
	return r.db.CreateTask(ctx, dbmodel.CreateTaskParams{
		UserID: userID,
		Title:  title,
		Description: sql.NullString{
			String: description,
			Valid:  description != "",
		},
		Priority: priority,
		DueAt: sql.NullTime{
			Time:  dueAt,
			Valid: !dueAt.IsZero(),
		},
	})
}

func (r *Repo) GetUserTasksRepo(
	ctx context.Context,
	userID uuid.UUID,
) ([]dbmodel.Task, error) {
	return r.db.GetUserTasks(ctx, userID)
}

func (r *Repo) GetTaskByIdRepo(
	ctx context.Context,
	taskID uuid.UUID,
	userID uuid.UUID,
) (dbmodel.Task, error) {
	return r.db.GetTaskById(ctx, dbmodel.GetTaskByIdParams{
		ID:     taskID,
		UserID: userID,
	})
}

func (r *Repo) UpdateTaskRepo(
	ctx context.Context,
	taskID uuid.UUID,
	userID uuid.UUID,
	title string,
	description string,
	priority string,
	status string,
	dueAt time.Time,
) (dbmodel.Task, error) {
	return r.db.UpdateTask(ctx, dbmodel.UpdateTaskParams{
		ID:     taskID,
		UserID: userID,
		Title:  title,
		Description: sql.NullString{
			String: description,
			Valid:  description != "",
		},
		Priority: priority,
		Status:   status,
		DueAt: sql.NullTime{
			Time:  dueAt,
			Valid: !dueAt.IsZero(),
		},
	})
}

func (r *Repo) DeleteTaskRepo(
	ctx context.Context,
	taskID uuid.UUID,
	userID uuid.UUID,
) error {
	return r.db.DeleteTask(ctx, dbmodel.DeleteTaskParams{
		ID:     taskID,
		UserID: userID,
	})
}
