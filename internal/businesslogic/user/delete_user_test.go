package user

import (
	"context"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/events"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestLogic_DeleteUser_StorageError(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	userId := "user-id"

	ts.mockUserStorage.EXPECT().DeleteUser(gomock.Any(), userId).
		Return(common.NewError(nil, common.ErrTypeNotFound))

	err := ts.userManager.DeleteUser(context.Background(), userId)
	assert.Error(t, err)
	var errCommon common.Error
	assert.ErrorAs(t, err, &errCommon)
	assert.Equal(t, common.ErrTypeNotFound, errCommon.Type())
}

func TestLogic_DeleteUser_EventEmitterError(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	userId := "user-id"

	ts.mockUserStorage.EXPECT().DeleteUser(gomock.Any(), userId).Return(nil)

	now := time.Now()
	ts.mockTimeProvider.EXPECT().Now().Return(now)

	userEvent := events.UserEvent{
		UserId:    userId,
		EventType: events.EventTypeDeleted,
		EventTime: now,
	}

	ts.mockEventEmitter.EXPECT().EmitUserEvent(gomock.Any(), userEvent).
		Return(common.NewError(nil, common.ErrTypeInternal))

	err := ts.userManager.DeleteUser(context.Background(), userId)
	assert.NoError(t, err) // Deletion should succeed even if event emission fails
}

func TestLogic_DeleteUser_Success(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	userId := "user-id"

	ts.mockUserStorage.EXPECT().DeleteUser(gomock.Any(), userId).Return(nil)

	now := time.Now()
	ts.mockTimeProvider.EXPECT().Now().Return(now)

	userEvent := events.UserEvent{
		UserId:    userId,
		EventType: events.EventTypeDeleted,
		EventTime: now,
	}

	ts.mockEventEmitter.EXPECT().EmitUserEvent(gomock.Any(), userEvent).Return(nil)

	err := ts.userManager.DeleteUser(context.Background(), userId)
	assert.NoError(t, err)
}
