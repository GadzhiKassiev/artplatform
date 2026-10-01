package activity

import (
	"context"
	"encoding/json"

	"artplatform/backend/internal/logging"
	"artplatform/backend/internal/service/activity/model"
	"artplatform/backend/internal/service/activity/storage"
	"artplatform/backend/internal/transport/kafka"
)

const (
	TopicCourseViewed    = "course.viewed"
	TopicPurchaseCreated = "purchase.created"
	ConsumerGroup        = "activity-worker"
)

type Service struct {
	storage storage.Storage
}

func New(st storage.Storage) *Service {
	return &Service{storage: st}
}

func (s *Service) GetStats(ctx context.Context, courseID string) (model.CourseStats, error) {
	return s.storage.GetStats(ctx, courseID)
}

func (s *Service) RecordView(ctx context.Context, courseID string) error {
	return s.storage.IncrementViews(ctx, courseID)
}

type CourseViewedEvent struct {
	CourseID string `json:"course_id"`
	UserID   string `json:"user_id"`
}

type PurchaseCreatedEvent struct {
	PurchaseID string  `json:"purchase_id"`
	UserID     string  `json:"user_id"`
	CourseID   string  `json:"course_id"`
	Amount     float64 `json:"amount"`
}

func (s *Service) StartConsumers(ctx context.Context, brokers []string) error {
	viewConsumer := kafka.NewConsumer(brokers, TopicCourseViewed, ConsumerGroup, s.handleCourseViewed)
	purchaseConsumer := kafka.NewConsumer(brokers, TopicPurchaseCreated, ConsumerGroup, s.handlePurchaseCreated)

	errCh := make(chan error, 2)
	go func() { errCh <- viewConsumer.Run(ctx) }()
	go func() { errCh <- purchaseConsumer.Run(ctx) }()

	select {
	case <-ctx.Done():
		viewConsumer.Close()
		purchaseConsumer.Close()
		return nil
	case err := <-errCh:
		return err
	}
}

func (s *Service) handleCourseViewed(ctx context.Context, key string, value []byte) error {
	var event CourseViewedEvent
	if err := json.Unmarshal(value, &event); err != nil {
		logging.ContextErrorE(ctx, "failed to decode course.viewed", err)
		return err
	}

	logging.ContextInfo(ctx, "recording course view", logging.NewKV("courseID", event.CourseID))
	return s.storage.IncrementViews(ctx, event.CourseID)
}

func (s *Service) handlePurchaseCreated(ctx context.Context, key string, value []byte) error {
	var event PurchaseCreatedEvent
	if err := json.Unmarshal(value, &event); err != nil {
		logging.ContextErrorE(ctx, "failed to decode purchase.created", err)
		return err
	}

	logging.ContextInfo(ctx, "recording purchase", logging.NewKV("courseID", event.CourseID))
	return s.storage.IncrementPurchases(ctx, event.CourseID)
}
