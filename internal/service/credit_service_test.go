package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"bank-app/internal/config"
	"bank-app/internal/model"
)

type MockCreditRepository struct {
	mock.Mock
}

func (m *MockCreditRepository) Create(ctx context.Context, credit *model.Credit) error {
	args := m.Called(ctx, credit)
	return args.Error(0)
}

func (m *MockCreditRepository) GetByID(ctx context.Context, id int64) (*model.Credit, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Credit), args.Error(1)
}

func (m *MockCreditRepository) GetByUserID(ctx context.Context, userID int64) ([]*model.Credit, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Credit), args.Error(1)
}

func (m *MockCreditRepository) GetSchedule(ctx context.Context, creditID int64) ([]*model.PaymentSchedule, error) {
	args := m.Called(ctx, creditID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.PaymentSchedule), args.Error(1)
}

func (m *MockCreditRepository) Update(ctx context.Context, credit *model.Credit) error {
	args := m.Called(ctx, credit)
	return args.Error(0)
}

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) Create(ctx context.Context, account *model.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockAccountRepository) GetByID(ctx context.Context, id int64) (*model.Account, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Account), args.Error(1)
}

func (m *MockAccountRepository) GetByUserID(ctx context.Context, userID int64) ([]*model.Account, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Account), args.Error(1)
}

func (m *MockAccountRepository) Update(ctx context.Context, account *model.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func TestCreditService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("успешное создание кредита", func(t *testing.T) {
		// Подготовка
		mockCreditRepo := new(MockCreditRepository)
		mockAccountRepo := new(MockAccountRepository)
		cfg := &config.Config{}
		service := NewCreditService(mockCreditRepo, mockAccountRepo, cfg)

		userID := int64(1)
		accountID := int64(1)
		amount := 100000.0
		term := 12

		account := &model.Account{
			ID:      accountID,
			UserID:  userID,
			Balance: 0,
		}

		// Проверяем, что нет активных кредитов
		mockCreditRepo.On("GetByUserID", ctx, userID).Return([]*model.Credit{}, nil)
		mockAccountRepo.On("GetByID", ctx, accountID).Return(account, nil)
		mockCreditRepo.On("Create", ctx, mock.AnythingOfType("*model.Credit")).Return(nil)
		mockAccountRepo.On("Update", ctx, mock.AnythingOfType("*model.Account")).Return(nil)

		// Действие
		err := service.Create(ctx, userID, accountID, amount, term)

		// Проверка
		assert.NoError(t, err)
		mockCreditRepo.AssertExpectations(t)
		mockAccountRepo.AssertExpectations(t)

		// Проверяем, что был вызван Create с правильными параметрами
		if calls := mockCreditRepo.Calls; len(calls) > 0 {
			for _, call := range calls {
				if call.Method == "Create" {
					credit := call.Arguments[1].(*model.Credit)
					assert.Equal(t, userID, credit.UserID)
					assert.Equal(t, accountID, credit.AccountID)
					assert.Equal(t, amount, credit.Amount)
					assert.Equal(t, term, credit.Term)
					assert.Equal(t, "active", credit.Status)
					assert.Greater(t, credit.MonthlyPayment, 0.0)
				}
			}
		}
	})

	t.Run("превышен лимит кредитной нагрузки", func(t *testing.T) {
		// Подготовка
		mockCreditRepo := new(MockCreditRepository)
		mockAccountRepo := new(MockAccountRepository)
		cfg := &config.Config{}
		service := NewCreditService(mockCreditRepo, mockAccountRepo, cfg)

		userID := int64(1)
		accountID := int64(1)
		amount := 100000.0
		term := 12

		account := &model.Account{
			ID:      accountID,
			UserID:  userID,
			Balance: 0,
		}

		existingCredits := []*model.Credit{
			{
				UserID: userID,
				Amount: 1000000,
				Status: "active",
			},
		}

		// Сначала должна быть проверка существования счета
		mockAccountRepo.On("GetByID", ctx, accountID).Return(account, nil).Once()
		// Затем проверка кредитной нагрузки
		mockCreditRepo.On("GetByUserID", ctx, userID).Return(existingCredits, nil).Once()

		// Действие
		err := service.Create(ctx, userID, accountID, amount, term)

		// Проверка
		assert.Error(t, err)
		assert.Equal(t, "credit load limit exceeded", err.Error())
		mockCreditRepo.AssertExpectations(t)
		mockAccountRepo.AssertExpectations(t)
	})

	t.Run("счет не найден", func(t *testing.T) {
		// Подготовка
		mockCreditRepo := new(MockCreditRepository)
		mockAccountRepo := new(MockAccountRepository)
		cfg := &config.Config{}
		service := NewCreditService(mockCreditRepo, mockAccountRepo, cfg)

		userID := int64(1)
		accountID := int64(999)
		amount := 100000.0
		term := 12

		mockAccountRepo.On("GetByID", ctx, accountID).Return(nil, errors.New("account not found"))

		// Действие
		err := service.Create(ctx, userID, accountID, amount, term)

		// Проверка
		assert.Error(t, err)
		mockAccountRepo.AssertExpectations(t)
	})
}

func TestCreditService_calculateMonthlyPayment(t *testing.T) {
	mockCreditRepo := new(MockCreditRepository)
	mockAccountRepo := new(MockAccountRepository)
	cfg := &config.Config{}
	service := &CreditSvc{
		repo:     mockCreditRepo,
		accounts: mockAccountRepo,
		cfg:      cfg,
	}

	tests := []struct {
		name     string
		amount   float64
		term     int
		rate     float64
		expected float64
	}{
		{
			name:     "кредит 100000 на 12 месяцев под 12%",
			amount:   100000,
			term:     12,
			rate:     12,
			expected: 8884.88,
		},
		{
			name:     "кредит 500000 на 24 месяца под 12%",
			amount:   500000,
			term:     24,
			rate:     12,
			expected: 23536.74,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment := service.calculateMonthlyPayment(tt.amount, tt.term, tt.rate)
			assert.InDelta(t, tt.expected, payment, 0.01)
		})
	}
}
