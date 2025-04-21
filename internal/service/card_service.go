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

type CardSvc struct {
	repo repository.CardRepository
}

func NewCardService(repo repository.CardRepository) CardService {
	return &CardSvc{repo: repo}
}

func (s *CardSvc) Create(ctx context.Context, accountID int64) error {
	// Генерируем номер карты (в реальном приложении использовать более надежный алгоритм)
	rand.Seed(time.Now().UnixNano())
	cardNumber := fmt.Sprintf("4276%012d", rand.Int63n(1000000000000))

	// Генерируем CVV
	cvv := fmt.Sprintf("%03d", rand.Intn(1000))

	// Срок действия карты (4 года от текущей даты)
	expiryDate := time.Now().AddDate(4, 0, 0)

	card := &model.Card{
		AccountID:  accountID,
		Number:     cardNumber,
		CVV:        cvv,
		ExpiryDate: expiryDate,
		Status:     "active",
	}

	return s.repo.Create(ctx, card)
}

func (s *CardSvc) GetByID(ctx context.Context, id int64) (*model.Card, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CardSvc) GetByAccountID(ctx context.Context, accountID int64) ([]*model.Card, error) {
	return s.repo.GetByAccountID(ctx, accountID)
}

func (s *CardSvc) Block(ctx context.Context, id int64) error {
	card, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	card.Status = "blocked"
	return s.repo.Update(ctx, card)
}

func (s *CardSvc) ValidateCard(ctx context.Context, number, cvv string) error {
	// В реальном приложении здесь должна быть проверка по алгоритму Луна
	// и проверка CVV через безопасное хранилище
	return errors.New("not implemented")
}
