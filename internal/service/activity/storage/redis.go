package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"artplatform/backend/internal/service/activity/model"
)

const statsTTL = 24 * time.Hour

type RedisStorage struct {
	client *redis.Client
}

func NewRedis(client *redis.Client) *RedisStorage {
	return &RedisStorage{client: client}
}

func viewsKey(courseID string) string {
	return fmt.Sprintf("stats:course:%s:views", courseID)
}

func purchasesKey(courseID string) string {
	return fmt.Sprintf("stats:course:%s:purchases", courseID)
}

func (s *RedisStorage) IncrementViews(ctx context.Context, courseID string) error {
	key := viewsKey(courseID)
	pipe := s.client.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, statsTTL)
	_, err := pipe.Exec(ctx)
	return err
}

func (s *RedisStorage) IncrementPurchases(ctx context.Context, courseID string) error {
	key := purchasesKey(courseID)
	pipe := s.client.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, statsTTL)
	_, err := pipe.Exec(ctx)
	return err
}

func (s *RedisStorage) GetStats(ctx context.Context, courseID string) (model.CourseStats, error) {
	pipe := s.client.Pipeline()
	viewsCmd := pipe.Get(ctx, viewsKey(courseID))
	purchasesCmd := pipe.Get(ctx, purchasesKey(courseID))

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return model.CourseStats{}, err
	}

	views, _ := viewsCmd.Int64()
	purchases, _ := purchasesCmd.Int64()

	return model.CourseStats{
		CourseID:      courseID,
		ViewCount:     views,
		PurchaseCount: purchases,
		UpdatedAt:     time.Now().UTC(),
	}, nil
}
