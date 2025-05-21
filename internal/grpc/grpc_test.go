package grpc

import (
	"github.com/alenalato/users-service/internal/businesslogic"
	"go.uber.org/mock/gomock"
	"testing"
)

type testSuite struct {
	mockCtrl        *gomock.Controller
	mockConverter   *MockmodelConverter
	mockUserManager *businesslogic.MockUserManager
	usersServer     *UsersServer
}

func newTestSuite(t *testing.T) *testSuite {
	mockCtrl := gomock.NewController(t)
	mockModelConverter := NewMockmodelConverter(mockCtrl)
	mockUserManager := businesslogic.NewMockUserManager(mockCtrl)

	usersServer := NewUsersServer(mockUserManager)
	usersServer.converter = mockModelConverter

	return &testSuite{
		mockCtrl:        mockCtrl,
		mockConverter:   mockModelConverter,
		mockUserManager: mockUserManager,
		usersServer:     usersServer,
	}
}
