package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"bank-app/internal/model"
	"bank-app/internal/repository"
)

type AccountSvc struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) AccountService {
	return &AccountSvc{repo: repo}
}

func (s *AccountSvc) Create(ctx context.Context, userID int64) error {
	// Генерируем номер счета (в реальном приложении использовать более надежный алгоритм)
	rand.Seed(time.Now().UnixNano())
	accountNumber := fmt.Sprintf("4080%011d", rand.Int63n(100000000000))

	account := &model.Account{
		UserID:   userID,
		Number:   accountNumber,
		Balance:  0,
		Currency: "RUB",
	}

	return s.repo.Create(ctx, account)
}

func (s *AccountSvc) GetByID(ctx context.Context, id int64) (*model.Account, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AccountSvc) GetByUserID(ctx context.Context, userID int64) ([]*model.Account, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *AccountSvc) UpdateBalance(ctx context.Context, id int64, amount float64) error {
	account, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	newBalance := account.Balance + amount
	if newBalance < 0 {
		return errors.New("insufficient funds")
	}

	account.Balance = newBalance
	return s.repo.Update(ctx, account)
}
