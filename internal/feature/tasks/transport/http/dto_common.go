package tasks_transport_http

import (
	"time"

	"github.com/DenisAstrakhan/api-server/internal/core/domain"
)

type TaskDTOResponse struct {
	ID           int        `json:"id" example:"1"`
	Version      int        `json:"version" example:"2"`
	Title        string     `json:"title" example:"Сделать домашку"`
	Description  *string    `json:"description" example:"Сделать Геометрия"`
	Completed    bool       `json:"completed" example:"true"`
	CreatedAT    time.Time  `json:"created_at" example:"2026-09-25T16:03:38.661621Z"`
	CompletedAt  *time.Time `json:"completed_at" example:"2026-09-29T09:59:47.655618Z"`
	AuthorUserID int        `json:"author_user_id" example:"1"`
}

func taskDTOFromDomain(task domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAT:    task.CreatedAT,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}
}
func tasksDTOFromDomains(tasks []domain.Task) []TaskDTOResponse {
	tasksDTO := make([]TaskDTOResponse, len(tasks))
	for i, task := range tasks {
		tasksDTO[i] = taskDTOFromDomain(task)
	}
	return tasksDTO
}
