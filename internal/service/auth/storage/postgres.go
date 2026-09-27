package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"artplatform/backend/internal/service/auth/model"
)

var ErrUserNotFound = errors.New("user not found")

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgres(pool *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{pool: pool}
}

func (s *PostgresStorage) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := `
		INSERT INTO users (id, email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := s.pool.Exec(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *PostgresStorage) GetUserByID(ctx context.Context, id model.UserID) (model.User, error) {
	query := `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var user model.User
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, ErrUserNotFound
	}
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *PostgresStorage) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	query := `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var user model.User
	err := s.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, ErrUserNotFound
	}
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
