package statistics_service

import (
	"context"
	"time"

	"github.com/DenisAstrakhan/api-server/internal/core/domain"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetTask(ctx context.Context, userID *int, from *time.Time, to *time.Time) ([]domain.Task, error)
}

func NewStatisticsService(statisticsRepository StatisticsRepository) StatisticsService {
	return StatisticsService{
		statisticsRepository: statisticsRepository,
	}
}
