package service

import (
	"context"

	"bank-app/internal/repository"
)

type AnalyticsSvc struct {
	repo repository.AnalyticsRepository
}

func NewAnalyticsService(repo repository.AnalyticsRepository) AnalyticsService {
	return &AnalyticsSvc{repo: repo}
}

func (s *AnalyticsSvc) GetTransactionAnalytics(ctx context.Context, userID int64, period string) (map[string]float64, error) {
	// В реальном приложении здесь должен быть анализ транзакций по категориям
	return map[string]float64{
		"income":      0,
		"expenses":    0,
		"transfers":   0,
		"deposits":    0,
		"withdrawals": 0,
	}, nil
}

func (s *AnalyticsSvc) GetCreditLoad(ctx context.Context, userID int64) (float64, error) {
	return s.repo.GetCreditLoad(ctx, userID)
}

func (s *AnalyticsSvc) PredictBalance(ctx context.Context, accountID int64, days int) (float64, error) {
	return s.repo.PredictBalance(ctx, accountID, days)
}
