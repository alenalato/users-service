package kafka

import (
	"errors"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/logger"
	"github.com/segmentio/kafka-go"
)

type Config struct {
	Addresses []string
}

type EventEmitter struct {
	topicName string
	writer    *kafka.Writer
}

func (e *EventEmitter) Close() error {
	return e.writer.Close()
}

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
