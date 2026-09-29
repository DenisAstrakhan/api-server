package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/DenisAstrakhan/api-server/internal/core/domain"
	core_errors "github.com/DenisAstrakhan/api-server/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(ctx context.Context, userID *int, from *time.Time, to *time.Time) (domain.Statistics, error) {
	if from != nil && to != nil {
		//смотру находится ли to до from или to=from
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf("'to' must be after 'from': %w", core_errors.ErrInvalidArgument)
		}
	}
	tasks, err := s.statisticsRepository.GetTask(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("get tasks from repository: %w", err)
	}
	statistics := calcStatistics(tasks)
	return statistics, nil
}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.NewStatistics(0, 0, nil, nil)
	}
	tasksCreated := len(tasks)

	var tasaksCompleted int
	var totalCompletionDuration time.Duration
	for _, v := range tasks {
		if v.Completed {
			tasaksCompleted++
		}
		completionDuration := v.CompletionDuration()
		if completionDuration != nil {
			totalCompletionDuration += *completionDuration
		}
	}

	tasaksCompletedRate := float64(tasaksCompleted) / float64(tasksCreated) * 100

	var tasksAverageCompletionTime time.Duration
	if tasaksCompleted > 0 && totalCompletionDuration != 0 {
		tasksAverageCompletionTime = totalCompletionDuration / time.Duration(tasaksCompleted)
	}
	return domain.NewStatistics(tasksCreated, tasaksCompleted, &tasaksCompletedRate, &tasksAverageCompletionTime)
}
