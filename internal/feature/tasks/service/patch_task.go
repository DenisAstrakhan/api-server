package task_service

import (
	"context"
	"fmt"

	"github.com/DenisAstrakhan/api-server/internal/core/domain"
)

func (s *TaskService) PatchTask(ctx context.Context, taskID int, patch domain.TaskPatch) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, taskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task from repository: %w", err)
	}
	if err := task.ApplyPatch(patch); err != nil {
		return domain.Task{}, fmt.Errorf("apply task patch: %w", err)
	}
	patchedTask, err := s.tasksRepository.PatchTask(ctx, taskID, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("patch task: %w", err)
	}
	return patchedTask, nil
}
