package storage

import (
	"context"

	"artplatform/backend/internal/service/purchase/model"
)

type Storage interface {
	CreatePurchase(ctx context.Context, p model.Purchase) (model.Purchase, error)
	GetPurchaseByID(ctx context.Context, id model.PurchaseID) (model.Purchase, error)
	GetPurchasesByUser(ctx context.Context, userID string) ([]model.Purchase, error)
	HasAccess(ctx context.Context, userID, courseID string) (bool, error)
	UpdatePurchase(ctx context.Context, p model.Purchase) (model.Purchase, error)
}
