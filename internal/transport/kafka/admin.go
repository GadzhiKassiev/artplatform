package kafka

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

// EnsureTopic creates the topic if it does not exist.
// Idempotent: safe to call multiple times.
func EnsureTopic(brokers []string, topic string) error {
	if len(brokers) == 0 {
		return fmt.Errorf("no brokers provided")
	}

	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return fmt.Errorf("dial kafka: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("get controller: %w", err)
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(
		controller.Host, strconv.Itoa(controller.Port),
	))
	if err != nil {
		return fmt.Errorf("dial controller: %w", err)
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	if err != nil {
		if err.Error() == "Topic with this name already exists." ||
			contains(err.Error(), "already exists") {
			return nil
		}
		return fmt.Errorf("create topic: %w", err)
	}

	return nil
}

// WaitForKafka retries EnsureTopic until it succeeds or timeout.
func WaitForKafka(brokers []string, topic string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	backoff := 500 * time.Millisecond

	for {
		err := EnsureTopic(brokers, topic)
		if err == nil {
			return nil
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("kafka not ready after %v: %w", timeout, err)
		}

		time.Sleep(backoff)
		if backoff < 5*time.Second {
			backoff *= 2
		}
	}
}

func contains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
