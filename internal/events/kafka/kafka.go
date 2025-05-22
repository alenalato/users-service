package kafka

import (
	"errors"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/logger"
	"github.com/segmentio/kafka-go"
)

// Config holds the configuration for the Kafka event emitter
type Config struct {
	Addresses []string
}

// EventEmitter is a struct that implements the EventEmitter interface using Kafka
type EventEmitter struct {
	// topicName is the name of the Kafka topic to which events will be emitted
	topicName string
	// writer is the Kafka writer used to send messages to the topic
	writer Writer
}

// Close closes the Kafka writer for proper resource cleanup
func (e *EventEmitter) Close() error {
	return e.writer.Close()
}

// NewEventEmitter creates a new EventEmitter instance
func NewEventEmitter(topicName string, config Config) (*EventEmitter, error) {
	if config.Addresses == nil {
		err := errors.New("addresses are empty")
		logger.Log.Error(err)

		return nil, common.NewError(err, common.ErrTypeInternal)
	}

	logger.Log.Debugf("Creating new kafka writer: %s@%v", topicName, config.Addresses)

	writer := &kafka.Writer{
		Addr:     kafka.TCP(config.Addresses...),
		Topic:    topicName,
		Balancer: &kafka.Hash{},
	}

	return &EventEmitter{
		topicName: topicName,
		writer:    writer,
	}, nil
}
