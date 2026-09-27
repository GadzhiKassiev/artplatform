package storage

import (
	"context"

	"artplatform/backend/internal/service/auth/model"
)

type Storage interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	GetUserByID(ctx context.Context, id model.UserID) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}
