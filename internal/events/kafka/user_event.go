package kafka

import (
	"context"
	"encoding/json"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/events"
	"github.com/alenalato/users-service/internal/logger"
	"github.com/segmentio/kafka-go"
)

// EmitUserEvent emits a user event to the Kafka topic
// It marshals the user event to JSON and sends it as a message
// to the Kafka topic using the Kafka writer
func (e *EventEmitter) EmitUserEvent(ctx context.Context, userEvent events.UserEvent) error {
	// Marshal user event to JSON
	userEventBytes, err := json.Marshal(userEvent)
	if err != nil {
		logger.Log.Errorf("Failed to marshal user event: %v", err)
		return err
	}

	if writeErr := e.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(userEvent.UserId),
		Value: userEventBytes,
	}); writeErr != nil {
		logger.Log.Errorf("Failed to emit user event: %v", writeErr)

		return common.NewError(writeErr, common.ErrTypeInternal)
	}

	return nil
}
