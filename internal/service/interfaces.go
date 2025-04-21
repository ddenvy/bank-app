package service

import (
	"bank-app/internal/model"
	"context"
)

type UserService interface {
	Register(ctx context.Context, username, email, password string) error
	Login(ctx context.Context, email, password string) (string, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
}

type AccountService interface {
	Create(ctx context.Context, userID int64) error
	GetByID(ctx context.Context, id int64) (*model.Account, error)
	GetByUserID(ctx context.Context, userID int64) ([]*model.Account, error)
	UpdateBalance(ctx context.Context, id int64, amount float64) error
}

type CardService interface {
	Create(ctx context.Context, accountID int64) error
	GetByID(ctx context.Context, id int64) (*model.Card, error)
	GetByAccountID(ctx context.Context, accountID int64) ([]*model.Card, error)
	ValidateCard(ctx context.Context, number, cvv string) error
}

type TransferService interface {
	Transfer(ctx context.Context, fromID, toID int64, amount float64) error
	GetByID(ctx context.Context, id int64) (*model.Transaction, error)
	GetByAccountID(ctx context.Context, accountID int64) ([]*model.Transaction, error)
}

type CreditService interface {
	Create(ctx context.Context, userID, accountID int64, amount float64, term int) error
	GetByID(ctx context.Context, id int64) (*model.Credit, error)
	GetSchedule(ctx context.Context, creditID int64) ([]*model.PaymentSchedule, error)
	ProcessPayments(ctx context.Context) error
}

type AnalyticsService interface {
	GetTransactionAnalytics(ctx context.Context, userID int64, period string) (map[string]float64, error)
	GetCreditLoad(ctx context.Context, userID int64) (float64, error)
	PredictBalance(ctx context.Context, accountID int64, days int) (float64, error)
}
