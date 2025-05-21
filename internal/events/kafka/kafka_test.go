package kafka

import (
	"go.uber.org/mock/gomock"
	"testing"
)

type testSuite struct {
	mockCtrl     *gomock.Controller
	mockWriter   *MockWriter
	eventEmitter *EventEmitter
}

func newTestSuite(t *testing.T) *testSuite {
	mockCtrl := gomock.NewController(t)
	mockWriter := NewMockWriter(mockCtrl)

	eventEmitter := &EventEmitter{
		topicName: "test-topic",
		writer:    mockWriter,
	}

	return &testSuite{
		mockCtrl:     mockCtrl,
		mockWriter:   mockWriter,
		eventEmitter: eventEmitter,
	}
}
