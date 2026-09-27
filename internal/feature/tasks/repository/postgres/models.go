package tasks_postgres_repository

import (
	"time"

	"github.com/DenisAstrakhan/api-server/internal/core/domain"
)

type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAT    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func tasksDomainsFromModels(tasks []TaskModel) []domain.Task {
	domainTasks := make([]domain.Task, len(tasks))
	for i, v := range tasks {
		domainTasks[i] = domain.NewTask(v.ID, v.Version, v.Title, v.Description, v.Completed, v.CreatedAT, v.CompletedAt, v.AuthorUserID)
	}
	return domainTasks
}
