package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"bank-app/internal/model"
)

// MockUserRepository - мок для репозитория пользователей
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func TestUserService_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("успешная регистрация", func(t *testing.T) {
		// Подготовка
		username := "testuser"
		email := "test@example.com"
		password := "password123"

		mockRepo.On("GetByEmail", ctx, email).Return(nil, errors.New("not found"))
		mockRepo.On("GetByUsername", ctx, username).Return(nil, errors.New("not found"))
		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(nil)

		// Действие
		err := service.Register(ctx, username, email, password)

		// Проверка
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("email уже существует", func(t *testing.T) {
		// Подготовка
		username := "testuser"
		email := "existing@example.com"
		password := "password123"

		existingUser := &model.User{
			ID:        1,
			Username:  "existing",
			Email:     email,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("GetByEmail", ctx, email).Return(existingUser, nil)

		// Действие
		err := service.Register(ctx, username, email, password)

		// Проверка
		assert.Error(t, err)
		assert.Equal(t, "user with this email already exists", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("успешный вход", func(t *testing.T) {
		// Подготовка
		email := "test@example.com"
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &model.User{
			ID:           1,
			Email:        email,
			PasswordHash: string(hashedPassword),
		}

		mockRepo.On("GetByEmail", ctx, email).Return(user, nil)

		// Действие
		token, err := service.Login(ctx, email, password)

		// Проверка
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("неверный пароль", func(t *testing.T) {
		// Подготовка
		email := "test@example.com"
		password := "wrongpassword"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

		user := &model.User{
			ID:           1,
			Email:        email,
			PasswordHash: string(hashedPassword),
		}

		mockRepo.On("GetByEmail", ctx, email).Return(user, nil)

		// Действие
		token, err := service.Login(ctx, email, password)

		// Проверка
		assert.Error(t, err)
		assert.Equal(t, "invalid email or password", err.Error())
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("пользователь найден", func(t *testing.T) {
		// Подготовка
		userID := int64(1)
		expectedUser := &model.User{
			ID:        userID,
			Username:  "testuser",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("GetByID", ctx, userID).Return(expectedUser, nil)

		// Действие
		user, err := service.GetByID(ctx, userID)

		// Проверка
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("пользователь не найден", func(t *testing.T) {
		// Подготовка
		userID := int64(999)
		mockRepo.On("GetByID", ctx, userID).Return(nil, errors.New("user not found"))

		// Действие
		user, err := service.GetByID(ctx, userID)

		// Проверка
		assert.Error(t, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}
