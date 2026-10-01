package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"artplatform/backend/internal/service/purchase/model"
)

var (
	ErrPurchaseNotFound = errors.New("purchase not found")
	ErrAlreadyPurchased = errors.New("course already purchased")
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgres(pool *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{pool: pool}
}

func (s *PostgresStorage) CreatePurchase(ctx context.Context, p model.Purchase) (model.Purchase, error) {
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now

	query := `
		INSERT INTO purchases (id, user_id, course_id, amount, status, transaction_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := s.pool.Exec(ctx, query,
		p.ID, p.UserID, p.CourseID, p.Amount, p.Status,
		nullString(p.TransactionID), p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return model.Purchase{}, ErrAlreadyPurchased
		}
		return model.Purchase{}, err
	}
	return p, nil
}

func (s *PostgresStorage) GetPurchaseByID(ctx context.Context, id model.PurchaseID) (model.Purchase, error) {
	query := `
		SELECT id, user_id, course_id, amount, status, COALESCE(transaction_id::text, ''), created_at, updated_at
		FROM purchases
		WHERE id = $1
	`
	var p model.Purchase
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.UserID, &p.CourseID, &p.Amount, &p.Status,
		&p.TransactionID, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Purchase{}, ErrPurchaseNotFound
	}
	if err != nil {
		return model.Purchase{}, err
	}
	return p, nil
}

func (s *PostgresStorage) GetPurchasesByUser(ctx context.Context, userID string) ([]model.Purchase, error) {
	query := `
		SELECT id, user_id, course_id, amount, status, COALESCE(transaction_id::text, ''), created_at, updated_at
		FROM purchases
		WHERE user_id = $1 AND status = 'COMPLETED'
		ORDER BY created_at DESC
	`
	rows, err := s.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []model.Purchase
	for rows.Next() {
		var p model.Purchase
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.CourseID, &p.Amount, &p.Status,
			&p.TransactionID, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		purchases = append(purchases, p)
	}
	return purchases, rows.Err()
}

func (s *PostgresStorage) HasAccess(ctx context.Context, userID, courseID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM purchases
			WHERE user_id = $1 AND course_id = $2 AND status = 'COMPLETED'
		)
	`
	var exists bool
	err := s.pool.QueryRow(ctx, query, userID, courseID).Scan(&exists)
	return exists, err
}

func (s *PostgresStorage) UpdatePurchase(ctx context.Context, p model.Purchase) (model.Purchase, error) {
	p.UpdatedAt = time.Now().UTC()

	query := `
		UPDATE purchases
		SET status = $2, transaction_id = $3, updated_at = $4
		WHERE id = $1
	`
	tag, err := s.pool.Exec(ctx, query,
		p.ID, p.Status, nullString(p.TransactionID), p.UpdatedAt,
	)
	if err != nil {
		return model.Purchase{}, err
	}
	if tag.RowsAffected() == 0 {
		return model.Purchase{}, ErrPurchaseNotFound
	}
	return p, nil
}

func nullString(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func isUniqueViolation(err error) bool {
	return err != nil && (contains(err.Error(), "duplicate key") || contains(err.Error(), "unique constraint"))
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
