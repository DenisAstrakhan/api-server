package task_service

import (
	"context"

	"github.com/DenisAstrakhan/api-server/internal/core/domain"
)

type TaskService struct {
	tasksRepository TasksRepository
}

type TasksRepository interface {
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	GetTsaks(ctx context.Context, userID *int, limit *int, offset *int) ([]domain.Task, error)
}

func NewTasksService(tasksRepository TasksRepository) TaskService {
	return TaskService{
		tasksRepository: tasksRepository,
	}
}
