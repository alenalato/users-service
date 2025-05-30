// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/alenalato/users-service/internal/events (interfaces: EventEmitter)
//
// Generated by this command:
//
//	mockgen -destination=events_mock.go -package=events github.com/alenalato/users-service/internal/events EventEmitter
//

// Package events is a generated GoMock package.
package events

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockEventEmitter is a mock of EventEmitter interface.
type MockEventEmitter struct {
	ctrl     *gomock.Controller
	recorder *MockEventEmitterMockRecorder
	isgomock struct{}
}

// MockEventEmitterMockRecorder is the mock recorder for MockEventEmitter.
type MockEventEmitterMockRecorder struct {
	mock *MockEventEmitter
}

// NewMockEventEmitter creates a new mock instance.
func NewMockEventEmitter(ctrl *gomock.Controller) *MockEventEmitter {
	mock := &MockEventEmitter{ctrl: ctrl}
	mock.recorder = &MockEventEmitterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventEmitter) EXPECT() *MockEventEmitterMockRecorder {
	return m.recorder
}

// EmitUserEvent mocks base method.
func (m *MockEventEmitter) EmitUserEvent(ctx context.Context, userEvent UserEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EmitUserEvent", ctx, userEvent)
	ret0, _ := ret[0].(error)
	return ret0
}

// EmitUserEvent indicates an expected call of EmitUserEvent.
func (mr *MockEventEmitterMockRecorder) EmitUserEvent(ctx, userEvent any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EmitUserEvent", reflect.TypeOf((*MockEventEmitter)(nil).EmitUserEvent), ctx, userEvent)
}
