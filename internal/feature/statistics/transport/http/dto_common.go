package statistics_transport_http

import (
	"github.com/DenisAstrakhan/api-server/internal/core/domain"
)

type StatisticsDTOResponse struct {
	TasksCreated               int      `json:"tasks_created" example:"2"`
	TasksCompleted             int      `json:"tasks_completed" example:"1"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate" example:"50"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time" example:"1m36.449448s"`
}

func statisticsDTOFromDomain(statistics domain.Statistics) StatisticsDTOResponse {
	var avgTime *string
	if statistics.TasksAverageCompletionTime != nil {
		duration := statistics.TasksAverageCompletionTime.String()
		avgTime = &duration
	}
	return StatisticsDTOResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TascsCompleted,
		TasksCompletedRate:         statistics.TascsCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}
