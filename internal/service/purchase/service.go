package purchase

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"artplatform/backend/internal/logging"
	purchaseerr "artplatform/backend/internal/service/purchase/err"
	"artplatform/backend/internal/service/purchase/model"
	"artplatform/backend/internal/service/purchase/storage"
	"artplatform/backend/internal/transport/kafka"
	paymentpb "artplatform/backend/proto/payment"
)

const TopicPurchaseCreated = "purchase.created"

type Service struct {
	storage       storage.Storage
	producer      *kafka.Producer
	paymentClient paymentpb.PaymentServiceClient
}

func New(st storage.Storage, producer *kafka.Producer, paymentClient paymentpb.PaymentServiceClient) *Service {
	return &Service{storage: st, producer: producer, paymentClient: paymentClient}
}

type CreateInput struct {
	UserID   string
	CourseID string
	Amount   float64
}

type PurchaseCreatedEvent struct {
	PurchaseID string  `json:"purchase_id"`
	UserID     string  `json:"user_id"`
	CourseID   string  `json:"course_id"`
	Amount     float64 `json:"amount"`
}

func (s *Service) Create(ctx context.Context, input CreateInput) (model.Purchase, error) {
	logging.ContextInfo(ctx, "creating purchase",
		logging.NewKV("userID", input.UserID),
		logging.NewKV("courseID", input.CourseID),
		logging.NewKV("amount", input.Amount),
	)

	hasAccess, err := s.storage.HasAccess(ctx, input.UserID, input.CourseID)
	if err != nil {
		logging.ContextErrorE(ctx, "failed to check access", err)
		return model.Purchase{}, err
	}
	if hasAccess {
		logging.ContextWarn(ctx, "already purchased", logging.NewKV("courseID", input.CourseID))
		return model.Purchase{}, purchaseerr.ErrAlreadyPurchased
	}

	payResp, err := s.paymentClient.ProcessPayment(ctx, &paymentpb.ProcessPaymentRequest{
		UserId:   input.UserID,
		CourseId: input.CourseID,
		Amount:   input.Amount,
	}, grpc.WaitForReady(true))
	if err != nil {
		logging.ContextErrorE(ctx, "payment failed", err)
		return model.Purchase{}, purchaseerr.ErrPaymentFailed
	}

	purchase := model.Purchase{
		ID:            uuid.NewString(),
		UserID:        input.UserID,
		CourseID:      input.CourseID,
		Amount:        input.Amount,
		Status:        model.StatusCompleted,
		TransactionID: &payResp.Transaction.Id,
	}

	created, err := s.storage.CreatePurchase(ctx, purchase)
	if err != nil {
		logging.ContextErrorE(ctx, "failed to create purchase", err)
		return model.Purchase{}, err
	}

	event := PurchaseCreatedEvent{
		PurchaseID: created.ID,
		UserID:     created.UserID,
		CourseID:   created.CourseID,
		Amount:     created.Amount,
	}
	if err := s.producer.Publish(ctx, TopicPurchaseCreated, created.ID, event); err != nil {
		logging.ContextErrorE(ctx, "failed to publish purchase.created", err)
	}

	logging.ContextInfo(ctx, "purchase created", logging.NewKV("purchaseID", created.ID))
	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id model.PurchaseID) (model.Purchase, error) {
	p, err := s.storage.GetPurchaseByID(ctx, id)
	if err == storage.ErrPurchaseNotFound {
		return model.Purchase{}, purchaseerr.ErrPurchaseNotFound
	}
	return p, err
}

func (s *Service) GetUserPurchases(ctx context.Context, userID string) ([]model.Purchase, error) {
	return s.storage.GetPurchasesByUser(ctx, userID)
}

func (s *Service) CheckAccess(ctx context.Context, userID, courseID string) (bool, error) {
	return s.storage.HasAccess(ctx, userID, courseID)
}
