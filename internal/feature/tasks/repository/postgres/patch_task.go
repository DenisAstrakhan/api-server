package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/DenisAstrakhan/api-server/internal/core/domain"
	core_errors "github.com/DenisAstrakhan/api-server/internal/core/errors"
	core_postgres_pool "github.com/DenisAstrakhan/api-server/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) PatchTask(ctx context.Context, taskId int, task domain.Task) (domain.Task, error) {
	timectx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	query := `
	UPDATE todoapp.tasks
	SET
		version=version+1,
		title=$1,
		description=$2,
		completed=$3,
		completed_at=$4		
	WHERE id=$5 AND version=$6
	RETURNING id, version, title, description, completed, created_at, completed_at, author_user_id;
	`
	row := r.pool.QueryRow(
		timectx,
		query,
		task.Title,
		task.Description,
		task.Completed,
		task.CompletedAt,
		taskId,
		task.Version,
	)
	var taskModel TaskModel
	err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreatedAT,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserID,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf("task with id='%d' concurrently accessed: %w", taskId, core_errors.ErrConflict)
		}
		return domain.Task{}, fmt.Errorf("scan task: %w", err)
	}

	return tasksDomainFromModel(taskModel), nil
}
