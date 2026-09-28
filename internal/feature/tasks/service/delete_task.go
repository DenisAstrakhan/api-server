package task_service

import (
	"context"
	"fmt"
)

func (s *TaskService) DeleteTask(ctx context.Context, taskID int) error {
	if err := s.tasksRepository.DeleteTask(ctx, taskID); err != nil {
		return fmt.Errorf("delet task: %w", err)
	}
	return nil
}
