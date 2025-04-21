package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"bank-app/internal/model"
	"bank-app/internal/repository"
)

type UserSvc struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserSvc{
		repo:      repo,
		jwtSecret: "your-secret-key", // В реальном приложении брать из конфига
	}
}

func (s *UserSvc) Register(ctx context.Context, username, email, password string) error {
	// Проверяем, что пользователь с таким email не существует
	if _, err := s.repo.GetByEmail(ctx, email); err == nil {
		return errors.New("user with this email already exists")
	}

	// Проверяем, что пользователь с таким username не существует
	if _, err := s.repo.GetByUsername(ctx, username); err == nil {
		return errors.New("user with this username already exists")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Создаем пользователя
	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	return s.repo.Create(ctx, user)
}

func (s *UserSvc) Login(ctx context.Context, email, password string) (string, error) {
	// Получаем пользователя по email
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Создаем JWT токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(user.ID, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	// Подписываем токен
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *UserSvc) GetByID(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}
