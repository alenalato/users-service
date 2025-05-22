package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

//go:generate mockgen -destination=writer_mock.go -package=kafka github.com/alenalato/users-service/internal/events/kafka Writer

// Writer is an interface for writing messages to Kafka
// It abstracts the underlying Kafka writer implementation
// to allow for easier testing and mocking
type Writer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}
