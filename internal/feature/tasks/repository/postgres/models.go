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
		domainTasks[i] = tasksDomainFromModel(v)
	}
	return domainTasks
}

func tasksDomainFromModel(model TaskModel) domain.Task {
	return domain.NewTask(
		model.ID,
		model.Version,
		model.Title,
		model.Description,
		model.Completed,
		model.CreatedAT,
		model.CompletedAt,
		model.AuthorUserID,
	)
}
