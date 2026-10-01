package storage

import (
	"context"

	"artplatform/backend/internal/service/payment/model"
)

type Storage interface {
	CreateTransaction(ctx context.Context, tx model.Transaction) (model.Transaction, error)
	GetTransactionByID(ctx context.Context, id model.TransactionID) (model.Transaction, error)
}
