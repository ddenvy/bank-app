package repository

import (
	"context"
	"database/sql"
	"errors"

	"bank-app/internal/model"
)

type TransferRepo struct {
	db *sql.DB
}

func NewTransferRepository(db *sql.DB) TransferRepository {
	return &TransferRepo{db: db}
}

func (r *TransferRepo) Create(ctx context.Context, transaction *model.Transaction) error {
	return errors.New("not implemented")
}

func (r *TransferRepo) GetByID(ctx context.Context, id int64) (*model.Transaction, error) {
	return nil, errors.New("not implemented")
}

func (r *TransferRepo) GetByAccountID(ctx context.Context, accountID int64) ([]*model.Transaction, error) {
	return nil, errors.New("not implemented")
}

func (r *TransferRepo) Update(ctx context.Context, transaction *model.Transaction) error {
	return errors.New("not implemented")
}
