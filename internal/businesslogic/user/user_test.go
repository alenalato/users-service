package user

import (
	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/events"
	"github.com/alenalato/users-service/internal/storage"
	"go.uber.org/mock/gomock"
	"testing"
)

type testSuite struct {
	mockCtrl            *gomock.Controller
	mockTimeProvider    *common.MockTimeProvider
	mockModelConverter  *MockmodelConverter
	mockPasswordManager *businesslogic.MockPasswordManager
	mockUserStorage     *storage.MockUserStorage
	mockEventEmitter    *events.MockEventEmitter
	userManager         *Logic
}

func newTestSuite(t *testing.T) *testSuite {
	mockCtrl := gomock.NewController(t)

	mockTimeProvider := common.NewMockTimeProvider(mockCtrl)
	mockModelConverter := NewMockmodelConverter(mockCtrl)
	mockPasswordManager := businesslogic.NewMockPasswordManager(mockCtrl)
	mockUserStorage := storage.NewMockUserStorage(mockCtrl)
	mockEventEmitter := events.NewMockEventEmitter(mockCtrl)

	userManager := NewLogic(
		mockPasswordManager,
		mockUserStorage,
		mockEventEmitter,
	)
	userManager.time = mockTimeProvider
	userManager.converter = mockModelConverter

	return &testSuite{
		mockCtrl:            mockCtrl,
		mockTimeProvider:    mockTimeProvider,
		mockModelConverter:  mockModelConverter,
		mockPasswordManager: mockPasswordManager,
		mockUserStorage:     mockUserStorage,
		mockEventEmitter:    mockEventEmitter,
		userManager:         userManager,
	}
}
