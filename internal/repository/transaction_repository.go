package repository

import (
	"context"
	"fmt"

	"github.com/azizanonerm/fintrack-api/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	pool *pgxpool.Pool
}

func NewTransactionRepository(pool *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{
		pool: pool,
	}
}

func (r *TransactionRepository) Create(
	ctx context.Context,
	t model.Transaction,
) (model.Transaction, error) {
	query := `
		INSERT INTO transactions (
			user_id,
			category_id,
			type,
			amount,
			description,
			transaction_date
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		t.UserID,
		t.CategoryID,
		t.Type,
		t.Amount,
		t.Description,
		t.TransactionDate,
	).Scan(
		&t.ID,
		&t.CreatedAt,
	)

	if err != nil {
		return model.Transaction{}, fmt.Errorf("create transaction: %w", err)
	}

	return t, nil
}

func (r *TransactionRepository) List(
	ctx context.Context,
	userID uuid.UUID,
) ([]model.Transaction, error) {
	query := `
		SELECT
			id,
			user_id,
			category_id,
			type,
			amount,
			description,
			transaction_date,
			created_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY transaction_date DESC, created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list transactions: %w", err)
	}
	defer rows.Close()

	var transactions []model.Transaction

	for rows.Next() {
		var transaction model.Transaction

		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.CategoryID,
			&transaction.Type,
			&transaction.Amount,
			&transaction.Description,
			&transaction.TransactionDate,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate transactions: %w", err)
	}

	return transactions, nil
}

func (r *TransactionRepository) Update(
	ctx context.Context,
	userID uuid.UUID,
	transactionID uuid.UUID,
	t model.Transaction,
) (model.Transaction, error) {
	query := `
		UPDATE transactions
		SET
			category_id = $1,
			type = $2,
			amount = $3,
			description = $4,
			transaction_date = $5
		WHERE id = $6 AND user_id = $7
		RETURNING
			id,
			user_id,
			category_id,
			type,
			amount,
			description,
			transaction_date,
			created_at
	`

	var updated model.Transaction

	err := r.pool.QueryRow(
		ctx,
		query,
		t.CategoryID,
		t.Type,
		t.Amount,
		t.Description,
		t.TransactionDate,
		transactionID,
		userID,
	).Scan(
		&updated.ID,
		&updated.UserID,
		&updated.CategoryID,
		&updated.Type,
		&updated.Amount,
		&updated.Description,
		&updated.TransactionDate,
		&updated.CreatedAt,
	)

	if err != nil {
		return model.Transaction{}, fmt.Errorf("update transaction: %w", err)
	}

	return updated, nil
}

func (r *TransactionRepository) Delete(
	ctx context.Context,
	userID uuid.UUID,
	transactionID uuid.UUID,
) error {
	result, err := r.pool.Exec(
		ctx,
		`DELETE FROM transactions WHERE id = $1 AND user_id = $2`,
		transactionID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("delete transaction: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
