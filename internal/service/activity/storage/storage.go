package storage

import (
	"context"

	"artplatform/backend/internal/service/activity/model"
)

type Storage interface {
	IncrementViews(ctx context.Context, courseID string) error
	IncrementPurchases(ctx context.Context, courseID string) error
	GetStats(ctx context.Context, courseID string) (model.CourseStats, error)
}
