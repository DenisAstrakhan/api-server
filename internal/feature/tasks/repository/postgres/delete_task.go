package tasks_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/DenisAstrakhan/api-server/internal/core/errors"
)

func (r *TasksRepository) DeleteTask(ctx context.Context, taskID int) error {
	timectx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	query := `
	DELETE FROM todoapp.tasks
	WHERE id=$1;
	`
	cmdTag, err := r.pool.Exec(timectx, query, taskID)
	if err != nil {
		return fmt.Errorf("exec query; %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		//запрос прошёл успешно но не на какие строки в таблици он не повлиял
		return fmt.Errorf("task with id='%d: %w", taskID, core_errors.ErrNotFound)
	}
	return nil
}
