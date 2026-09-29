package domain

import "time"

type Statistics struct {
	TasksCreated               int
	TascsCompleted             int
	TascsCompletedRate         *float64
	TasksAverageCompletionTime *time.Duration
}

func NewStatistics(
	tasksCreated int,
	tascsCompleted int,
	tascsCompletedRate *float64,
	tasksAverageCompletionTime *time.Duration,
) Statistics {
	return Statistics{
		TasksCreated:               tasksCreated,
		TascsCompleted:             tascsCompleted,
		TascsCompletedRate:         tascsCompletedRate,
		TasksAverageCompletionTime: tasksAverageCompletionTime,
	}
}
