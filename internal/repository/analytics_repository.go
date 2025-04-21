package repository

import (
	"context"
	"database/sql"
	"errors"

	"bank-app/internal/model"
)

type AnalyticsRepo struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) AnalyticsRepository {
	return &AnalyticsRepo{db: db}
}

func (r *AnalyticsRepo) GetTransactionsByPeriod(ctx context.Context, userID int64, from, to string) ([]*model.Transaction, error) {
	return nil, errors.New("not implemented")
}

func (r *AnalyticsRepo) GetCreditLoad(ctx context.Context, userID int64) (float64, error) {
	return 0, errors.New("not implemented")
}

func (r *AnalyticsRepo) PredictBalance(ctx context.Context, accountID int64, days int) (float64, error) {
	return 0, errors.New("not implemented")
}
