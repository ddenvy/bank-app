package service

import (
	"context"
	"errors"

	"bank-app/internal/model"
	"bank-app/internal/repository"
)

type TransferSvc struct {
	repo     repository.TransferRepository
	accounts repository.AccountRepository
}

func NewTransferService(repo repository.TransferRepository, accounts repository.AccountRepository) TransferService {
	return &TransferSvc{
		repo:     repo,
		accounts: accounts,
	}
}

func (s *TransferSvc) Transfer(ctx context.Context, fromID, toID int64, amount float64) error {
	// Проверяем существование счетов
	fromAcc, err := s.accounts.GetByID(ctx, fromID)
	if err != nil {
		return err
	}

	toAcc, err := s.accounts.GetByID(ctx, toID)
	if err != nil {
		return err
	}

	// Проверяем достаточность средств
	if fromAcc.Balance < amount {
		return errors.New("insufficient funds")
	}

	// Создаем транзакцию
	transaction := &model.Transaction{
		FromAccountID: fromID,
		ToAccountID:   toID,
		Amount:        amount,
		Status:        "pending",
	}

	// Сохраняем транзакцию
	if err := s.repo.Create(ctx, transaction); err != nil {
		return err
	}

	// Обновляем балансы счетов
	fromAcc.Balance -= amount
	toAcc.Balance += amount

	if err := s.accounts.Update(ctx, fromAcc); err != nil {
		return err
	}

	if err := s.accounts.Update(ctx, toAcc); err != nil {
		// В реальном приложении здесь нужно откатить изменения первого счета
		return err
	}

	// Обновляем статус транзакции
	transaction.Status = "completed"
	return s.repo.Update(ctx, transaction)
}

func (s *TransferSvc) GetByID(ctx context.Context, id int64) (*model.Transaction, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TransferSvc) GetByAccountID(ctx context.Context, accountID int64) ([]*model.Transaction, error) {
	return s.repo.GetByAccountID(ctx, accountID)
}
