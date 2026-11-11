package repository

import (
	"context"
	"fmt"

	"github.com/azizanonerm/fintrack-api/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct{ pool *pgxpool.Pool }

func NewTransactionRepository(pool *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{pool: pool}
}

func (r *TransactionRepository) Create(ctx context.Context, t model.Transaction) (model.Transaction, error) {
	query := `INSERT INTO transactions (user_id, category_id, type, amount, description, transaction_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`
	err := r.pool.QueryRow(ctx, query, t.UserID, t.CategoryID, t.Type, t.Amount, t.Description, t.TransactionDate).
		Scan(&t.ID, &t.CreatedAt)
	if err != nil {
		return model.Transaction{}, fmt.Errorf("create transaction: %w", err)
	}
	return t, nil
}

func (r *TransactionRepository) List(ctx context.Context, userID uuid.UUID) ([]model.Transaction, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, user_id, category_id, type, amount, description, transaction_date, created_at FROM transactions WHERE user_id = $1 ORDER BY transaction_date DESC, created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("list transactions: %w", err)
	}
	defer rows.Close()

	var result []model.Transaction
	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.CategoryID, &t.Type, &t.Amount, &t.Description, &t.TransactionDate, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}
		result = append(result, t)
	}
	return result, rows.Err()
}
