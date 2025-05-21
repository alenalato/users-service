package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/alenalato/users-service/internal/common"
	"go.uber.org/mock/gomock"
	"testing"
	"time"

	"github.com/alenalato/users-service/internal/events"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func TestEmitUserEvent_Success(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	userEvent := events.UserEvent{
		EventType: events.EventTypeCreated,
		EventTime: time.Now(),
		UserId:    "123",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	userEventBytes, _ := json.Marshal(userEvent)

	ts.mockWriter.EXPECT().WriteMessages(gomock.Any(), kafka.Message{
		Key:   []byte(userEvent.UserId),
		Value: userEventBytes,
	}).Return(nil)

	err := ts.eventEmitter.EmitUserEvent(context.Background(), userEvent)

	assert.NoError(t, err)
}

func TestEmitUserEvent_WriteError(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	userEvent := events.UserEvent{
		EventType: events.EventTypeCreated,
		EventTime: time.Now(),
		UserId:    "123",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	userEventBytes, _ := json.Marshal(userEvent)

	ts.mockWriter.EXPECT().WriteMessages(gomock.Any(), kafka.Message{
		Key:   []byte(userEvent.UserId),
		Value: userEventBytes,
	}).Return(errors.New("write error"))

	err := ts.eventEmitter.EmitUserEvent(context.Background(), userEvent)

	assert.Error(t, err)
	var errCommon common.Error
	assert.ErrorAs(t, err, &errCommon)
	assert.Equal(t, common.ErrTypeInternal, errCommon.Type())
}
