package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/azizanonerm/fintrack-api/internal/model"
	"github.com/azizanonerm/fintrack-api/internal/repository"
	"github.com/google/uuid"
)

type TransactionService struct {
	repository *repository.TransactionRepository
}

type CreateTransactionInput struct {
	CategoryID      uuid.UUID
	Type            string
	Amount          float64
	Description     string
	TransactionDate time.Time
}

func NewTransactionService(r *repository.TransactionRepository) *TransactionService {
	return &TransactionService{repository: r}
}

func (s *TransactionService) Create(ctx context.Context, userID uuid.UUID, input CreateTransactionInput) (model.Transaction, error) {
	transactionType := strings.ToLower(strings.TrimSpace(input.Type))
	if transactionType != "income" && transactionType != "expense" {
		return model.Transaction{}, errors.New("type must be income or expense")
	}
	if input.Amount <= 0 {
		return model.Transaction{}, errors.New("amount must be greater than zero")
	}
	if input.CategoryID == uuid.Nil {
		return model.Transaction{}, errors.New("category_id is required")
	}
	if input.TransactionDate.IsZero() {
		return model.Transaction{}, errors.New("transaction_date is required")
	}

	return s.repository.Create(ctx, model.Transaction{UserID: userID, CategoryID: input.CategoryID, Type: transactionType, Amount: input.Amount, Description: strings.TrimSpace(input.Description), TransactionDate: input.TransactionDate})
}

func (s *TransactionService) List(ctx context.Context, userID uuid.UUID) ([]model.Transaction, error) {
	return s.repository.List(ctx, userID)
}

func (s *TransactionService) Update(ctx context.Context, userID, transactionID uuid.UUID, input CreateTransactionInput) (model.Transaction, error) {
	transactionType := strings.ToLower(strings.TrimSpace(input.Type))
	if transactionType != "income" && transactionType != "expense" {
		return model.Transaction{}, errors.New("type must be income or expense")
	}
	if input.Amount <= 0 {
		return model.Transaction{}, errors.New("amount must be greater than zero")
	}
	if input.CategoryID == uuid.Nil {
		return model.Transaction{}, errors.New("category_id is required")
	}
	if input.TransactionDate.IsZero() {
		return model.Transaction{}, errors.New("transaction_date is required")
	}

	return s.repository.Update(ctx, userID, transactionID, model.Transaction{
		UserID:          userID,
		CategoryID:      input.CategoryID,
		Type:            transactionType,
		Amount:          input.Amount,
		Description:     strings.TrimSpace(input.Description),
		TransactionDate: input.TransactionDate,
	})
}

func (s *TransactionService) Delete(ctx context.Context, userID, transactionID uuid.UUID) error {
	return s.repository.Delete(ctx, userID, transactionID)
}
