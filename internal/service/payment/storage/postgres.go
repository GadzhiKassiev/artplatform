package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"artplatform/backend/internal/service/payment/model"
)

var ErrTransactionNotFound = errors.New("transaction not found")

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgres(pool *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{pool: pool}
}

func (s *PostgresStorage) CreateTransaction(ctx context.Context, tx model.Transaction) (model.Transaction, error) {
	now := time.Now().UTC()
	tx.CreatedAt = now
	tx.UpdatedAt = now

	query := `
		INSERT INTO transactions (id, user_id, course_id, amount, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := s.pool.Exec(ctx, query,
		tx.ID, tx.UserID, tx.CourseID, tx.Amount, tx.Status, tx.CreatedAt, tx.UpdatedAt,
	)
	if err != nil {
		return model.Transaction{}, err
	}
	return tx, nil
}

func (s *PostgresStorage) GetTransactionByID(ctx context.Context, id model.TransactionID) (model.Transaction, error) {
	query := `
		SELECT id, user_id, course_id, amount, status, created_at, updated_at
		FROM transactions
		WHERE id = $1
	`
	var tx model.Transaction
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&tx.ID, &tx.UserID, &tx.CourseID, &tx.Amount, &tx.Status,
		&tx.CreatedAt, &tx.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Transaction{}, ErrTransactionNotFound
	}
	if err != nil {
		return model.Transaction{}, err
	}
	return tx, nil
}
