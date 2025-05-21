package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

//go:generate mockgen -destination=writer_mock.go -package=kafka github.com/alenalato/users-service/internal/events/kafka Writer

type Writer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}
