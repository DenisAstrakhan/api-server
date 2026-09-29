package statistics_postgres_repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/DenisAstrakhan/api-server/internal/core/domain"
)

func (r *StatisticsRepository) GetTask(ctx context.Context, userID *int, from *time.Time, to *time.Time) ([]domain.Task, error) {
	timectx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var queryBuilder strings.Builder

	queryBuilder.WriteString(`
		SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
		FROM todoapp.tasks
	`)
	args := []any{}
	condition := []string{}

	if userID != nil {
		condition = append(condition, fmt.Sprintf("author_user_id=$%d", len(args)+1))
		args = append(args, userID)
	}

	if from != nil {
		condition = append(condition, fmt.Sprintf("created_at>=$%d", len(args)+1))
		args = append(args, from)
	}

	if to != nil {
		condition = append(condition, fmt.Sprintf("created_at<$%d", len(args)+1))
		args = append(args, to)
	}

	if len(condition) > 0 {
		queryBuilder.WriteString(" WHERE " + strings.Join(condition, " AND "))
	}

	queryBuilder.WriteString(" ORDER BY id ASC")

	rows, err := r.pool.Query(timectx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()
	var tasksModels []TaskModel
	for rows.Next() {
		var taskModel TaskModel
		err := rows.Scan(
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
			return nil, fmt.Errorf("scan users: %w", err)
		}
		tasksModels = append(tasksModels, taskModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}
	tasksDomain := tasksDomainsFromModels(tasksModels)
	return tasksDomain, nil
}
