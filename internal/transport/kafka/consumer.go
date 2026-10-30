package kafka

import (
	"artplatform/backend/internal/logging"
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type MessageHandler func(ctx context.Context, key string, value []byte) error

type Consumer struct {
	reader  *kafka.Reader
	handler MessageHandler
}

func NewConsumer(brokers []string, topic, groupID string, handler MessageHandler) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 1,
			MaxBytes: 10e6,
		}),
		handler: handler,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	logging.ContextInfo(ctx, "kafka consumer started",
		logging.NewKV("topic", c.reader.Config().Topic),
		logging.NewKV("groupID", c.reader.Config().GroupID),
	)
	const (
		initialBackoff = 1 * time.Second
		maxBackoff     = 30 * time.Second
	)

	timer := time.NewTimer(initialBackoff)
	defer timer.Stop()

	backoff := initialBackoff

	for {
		if ctx.Err() != nil {
			return nil
		}
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}

			logging.ContextWarn(ctx, "kafka read error, retrying",
				logging.NewKV("error", err),
				logging.NewKV("backoff", backoff),
			)

			timer.Reset(backoff)
			select {
			case <-timer.C:
				backoff *= 2
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
				continue
			case <-ctx.Done():
				return nil
			}
		}

		backoff = initialBackoff
		taskCtx := context.Background()
		if err := c.handler(taskCtx, string(msg.Key), msg.Value); err != nil {
			logging.ContextErrorE(taskCtx, "handler failed", err)
			continue
		}
		if err := c.reader.CommitMessages(taskCtx, msg); err != nil {
			logging.ContextErrorE(taskCtx, "failed to commit offset", err)
			continue
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func Decode[T any](value []byte) (T, error) {
	var v T
	err := json.Unmarshal(value, &v)
	return v, err
}
