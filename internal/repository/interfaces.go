package repository

import (
	"bank-app/internal/model"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
}

type AccountRepository interface {
	Create(ctx context.Context, account *model.Account) error
	GetByID(ctx context.Context, id int64) (*model.Account, error)
	GetByUserID(ctx context.Context, userID int64) ([]*model.Account, error)
	Update(ctx context.Context, account *model.Account) error
}

type CardRepository interface {
	Create(ctx context.Context, card *model.Card) error
	GetByID(ctx context.Context, id int64) (*model.Card, error)
	GetByAccountID(ctx context.Context, accountID int64) ([]*model.Card, error)
	Update(ctx context.Context, card *model.Card) error
}

type TransferRepository interface {
	Create(ctx context.Context, transaction *model.Transaction) error
	GetByID(ctx context.Context, id int64) (*model.Transaction, error)
	GetByAccountID(ctx context.Context, accountID int64) ([]*model.Transaction, error)
	Update(ctx context.Context, transaction *model.Transaction) error
}

type CreditRepository interface {
	Create(ctx context.Context, credit *model.Credit) error
	GetByID(ctx context.Context, id int64) (*model.Credit, error)
	GetByUserID(ctx context.Context, userID int64) ([]*model.Credit, error)
	Update(ctx context.Context, credit *model.Credit) error
	GetSchedule(ctx context.Context, creditID int64) ([]*model.PaymentSchedule, error)
}

type AnalyticsRepository interface {
	GetTransactionsByPeriod(ctx context.Context, userID int64, from, to string) ([]*model.Transaction, error)
	GetCreditLoad(ctx context.Context, userID int64) (float64, error)
	PredictBalance(ctx context.Context, accountID int64, days int) (float64, error)
}
