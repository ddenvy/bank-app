package repository

import (
	"context"
	"database/sql"
	"errors"

	"bank-app/internal/model"
)

type CardRepo struct {
	db *sql.DB
}

func NewCardRepository(db *sql.DB) CardRepository {
	return &CardRepo{db: db}
}

func (r *CardRepo) Create(ctx context.Context, card *model.Card) error {
	return errors.New("not implemented")
}

func (r *CardRepo) GetByID(ctx context.Context, id int64) (*model.Card, error) {
	return nil, errors.New("not implemented")
}

func (r *CardRepo) GetByAccountID(ctx context.Context, accountID int64) ([]*model.Card, error) {
	return nil, errors.New("not implemented")
}

func (r *CardRepo) Update(ctx context.Context, card *model.Card) error {
	return errors.New("not implemented")
}
