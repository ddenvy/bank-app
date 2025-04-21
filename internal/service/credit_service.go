package service

import (
	"context"
	"errors"
	"math"

	"bank-app/internal/config"
	"bank-app/internal/model"
	"bank-app/internal/repository"
)

type CreditSvc struct {
	repo     repository.CreditRepository
	accounts repository.AccountRepository
	cfg      *config.Config
}

func NewCreditService(repo repository.CreditRepository, accounts repository.AccountRepository, cfg *config.Config) CreditService {
	return &CreditSvc{
		repo:     repo,
		accounts: accounts,
		cfg:      cfg,
	}
}

func (s *CreditSvc) Create(ctx context.Context, userID, accountID int64, amount float64, term int) error {
	// Проверяем существование счета
	account, err := s.accounts.GetByID(ctx, accountID)
	if err != nil {
		return err
	}

	// Получаем существующие кредиты
	credits, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Проверяем кредитную нагрузку с учетом нового кредита
	var totalAmount float64
	for _, credit := range credits {
		if credit.Status == "active" {
			totalAmount += credit.Amount
		}
	}

	// Проверяем общую сумму с учетом нового кредита
	if totalAmount+amount >= 1000000 {
		return errors.New("credit load limit exceeded")
	}

	// Рассчитываем ежемесячный платеж (12% годовых)
	monthlyPayment := s.calculateMonthlyPayment(amount, term, 12.0)

	credit := &model.Credit{
		UserID:         userID,
		AccountID:      accountID,
		Amount:         amount,
		Term:           term,
		Status:         "active",
		MonthlyPayment: monthlyPayment,
	}

	if err := s.repo.Create(ctx, credit); err != nil {
		return err
	}

	// Зачисляем сумму кредита на счет
	account.Balance += amount
	return s.accounts.Update(ctx, account)
}

func (s *CreditSvc) GetByID(ctx context.Context, id int64) (*model.Credit, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CreditSvc) GetByUserID(ctx context.Context, userID int64) ([]*model.Credit, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *CreditSvc) GetSchedule(ctx context.Context, creditID int64) ([]*model.PaymentSchedule, error) {
	return s.repo.GetSchedule(ctx, creditID)
}

func (s *CreditSvc) ProcessPayments(ctx context.Context) error {
	// В реальном приложении здесь должна быть обработка платежей по кредитам
	// с учетом графика платежей и баланса счетов
	return errors.New("not implemented")
}

func (s *CreditSvc) calculateMonthlyPayment(amount float64, term int, rate float64) float64 {
	monthlyRate := rate / 12 / 100
	payment := amount * monthlyRate * math.Pow(1+monthlyRate, float64(term)) /
		(math.Pow(1+monthlyRate, float64(term)) - 1)
	return math.Round(payment*100) / 100
}
