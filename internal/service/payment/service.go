package payment

import (
	"context"

	"github.com/google/uuid"

	"artplatform/backend/internal/logging"
	paymenterr "artplatform/backend/internal/service/payment/err"
	"artplatform/backend/internal/service/payment/model"
	"artplatform/backend/internal/service/payment/storage"
	"artplatform/backend/internal/transport/kafka"
)

const (
	TopicPaymentSucceeded = "payment.succeeded"
)

type Service struct {
	storage  storage.Storage
	producer *kafka.Producer
}

func New(st storage.Storage, producer *kafka.Producer) *Service {
	return &Service{storage: st, producer: producer}
}

type ProcessInput struct {
	UserID   string
	CourseID string
	Amount   float64
}

type PaymentSucceededEvent struct {
	TransactionID string  `json:"transaction_id"`
	UserID        string  `json:"user_id"`
	CourseID      string  `json:"course_id"`
	Amount        float64 `json:"amount"`
}

// ProcessPayment is a stub — always succeeds for prototype.
func (s *Service) ProcessPayment(ctx context.Context, input ProcessInput) (model.Transaction, error) {
	logging.ContextInfo(ctx, "processing payment",
		logging.NewKV("userID", input.UserID),
		logging.NewKV("courseID", input.CourseID),
		logging.NewKV("amount", input.Amount),
	)

	tx := model.Transaction{
		ID:       uuid.NewString(),
		UserID:   input.UserID,
		CourseID: input.CourseID,
		Amount:   input.Amount,
		Status:   model.StatusSuccess,
	}

	created, err := s.storage.CreateTransaction(ctx, tx)
	if err != nil {
		logging.ContextErrorE(ctx, "failed to save transaction", err)
		return model.Transaction{}, err
	}

	event := PaymentSucceededEvent{
		TransactionID: created.ID,
		UserID:        created.UserID,
		CourseID:      created.CourseID,
		Amount:        created.Amount,
	}
	if err := s.producer.Publish(ctx, TopicPaymentSucceeded, created.ID, event); err != nil {
		logging.ContextErrorE(ctx, "failed to publish payment.succeeded", err)
		// Not fatal — payment processed, event lost
	}

	logging.ContextInfo(ctx, "payment succeeded", logging.NewKV("transactionID", created.ID))
	return created, nil
}

func (s *Service) GetTransaction(ctx context.Context, id model.TransactionID) (model.Transaction, error) {
	tx, err := s.storage.GetTransactionByID(ctx, id)
	if err != nil {
		if err == storage.ErrTransactionNotFound {
			return model.Transaction{}, paymenterr.ErrTransactionNotFound
		}
		return model.Transaction{}, err
	}
	return tx, nil
}
