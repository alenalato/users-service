// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/alenalato/users-service/internal/storage (interfaces: UserStorage)
//
// Generated by this command:
//
//	mockgen -destination=storage_mock.go -package=storage github.com/alenalato/users-service/internal/storage UserStorage
//

// Package storage is a generated GoMock package.
package storage

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUserStorage is a mock of UserStorage interface.
type MockUserStorage struct {
	ctrl     *gomock.Controller
	recorder *MockUserStorageMockRecorder
	isgomock struct{}
}

// MockUserStorageMockRecorder is the mock recorder for MockUserStorage.
type MockUserStorageMockRecorder struct {
	mock *MockUserStorage
}

// NewMockUserStorage creates a new mock instance.
func NewMockUserStorage(ctrl *gomock.Controller) *MockUserStorage {
	mock := &MockUserStorage{ctrl: ctrl}
	mock.recorder = &MockUserStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserStorage) EXPECT() *MockUserStorageMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUserStorage) CreateUser(ctx context.Context, userDetails UserDetails) (*User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, userDetails)
	ret0, _ := ret[0].(*User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserStorageMockRecorder) CreateUser(ctx, userDetails any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserStorage)(nil).CreateUser), ctx, userDetails)
}

// DeleteUser mocks base method.
func (m *MockUserStorage) DeleteUser(ctx context.Context, userId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockUserStorageMockRecorder) DeleteUser(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserStorage)(nil).DeleteUser), ctx, userId)
}

// ListUsers mocks base method.
func (m *MockUserStorage) ListUsers(ctx context.Context, userFilter UserFilter, pageSize int, pageToken string) ([]User, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUsers", ctx, userFilter, pageSize, pageToken)
	ret0, _ := ret[0].([]User)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListUsers indicates an expected call of ListUsers.
func (mr *MockUserStorageMockRecorder) ListUsers(ctx, userFilter, pageSize, pageToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsers", reflect.TypeOf((*MockUserStorage)(nil).ListUsers), ctx, userFilter, pageSize, pageToken)
}

// UpdateUser mocks base method.
func (m *MockUserStorage) UpdateUser(ctx context.Context, userId string, userUpdate UserUpdate) (*User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, userId, userUpdate)
	ret0, _ := ret[0].(*User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserStorageMockRecorder) UpdateUser(ctx, userId, userUpdate any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserStorage)(nil).UpdateUser), ctx, userId, userUpdate)
}
