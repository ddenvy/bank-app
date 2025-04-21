package repository

import (
	"context"
	"database/sql"
	"errors"

	"bank-app/internal/model"
)

type CreditRepo struct {
	db *sql.DB
}

func NewCreditRepository(db *sql.DB) CreditRepository {
	return &CreditRepo{db: db}
}

func (r *CreditRepo) Create(ctx context.Context, credit *model.Credit) error {
	return errors.New("not implemented")
}

func (r *CreditRepo) GetByID(ctx context.Context, id int64) (*model.Credit, error) {
	return nil, errors.New("not implemented")
}

func (r *CreditRepo) GetByUserID(ctx context.Context, userID int64) ([]*model.Credit, error) {
	return nil, errors.New("not implemented")
}

func (r *CreditRepo) Update(ctx context.Context, credit *model.Credit) error {
	return errors.New("not implemented")
}

func (r *CreditRepo) GetSchedule(ctx context.Context, creditID int64) ([]*model.PaymentSchedule, error) {
	return nil, errors.New("not implemented")
}
