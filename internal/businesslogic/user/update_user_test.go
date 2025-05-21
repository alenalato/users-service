package user

import (
	"context"
	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/events"
	"github.com/alenalato/users-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestLogic_UpdateUser_ValidationError(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	userUpdate := businesslogic.UserUpdate{}

	res, err := ts.userManager.UpdateUser(context.Background(), "user-id", userUpdate)
	assert.Nil(t, res)
	var errCommon common.Error
	assert.ErrorAs(t, err, &errCommon)
	assert.Equal(t, common.ErrTypeInvalidArgument, errCommon.Type())
}

func TestLogic_UpdateUser_ConverterError(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	userUpdate := businesslogic.UserUpdate{
		FirstName:  "John",
		LastName:   "Doe",
		UpdateMask: []string{"first_name", "last_name"},
	}

	ts.mockModelConverter.EXPECT().fromModelUserUpdateToStorage(gomock.Any(), userUpdate).
		Return(storage.UserUpdate{}, common.NewError(nil, common.ErrTypeInternal))

	res, err := ts.userManager.UpdateUser(context.Background(), "user-id", userUpdate)
	assert.Nil(t, res)
	var errCommon common.Error
	assert.ErrorAs(t, err, &errCommon)
	assert.Equal(t, common.ErrTypeInternal, errCommon.Type())
}

func TestLogic_UpdateUser_StorageError(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	userUpdate := businesslogic.UserUpdate{
		FirstName:  "John",
		LastName:   "Doe",
		UpdateMask: []string{"first_name", "last_name"},
	}

	storageUserUpdate := storage.UserUpdate{
		FirstName: &userUpdate.FirstName,
		LastName:  &userUpdate.LastName,
	}

	ts.mockModelConverter.EXPECT().fromModelUserUpdateToStorage(gomock.Any(), userUpdate).Return(storageUserUpdate, nil)

	now := time.Now()
	ts.mockTimeProvider.EXPECT().Now().Return(now)

	storageUserUpdate.UpdatedAt = &now

	ts.mockUserStorage.EXPECT().UpdateUser(gomock.Any(), "user-id", storageUserUpdate).
		Return(nil, common.NewError(nil, common.ErrTypeNotFound))

	res, err := ts.userManager.UpdateUser(context.Background(), "user-id", userUpdate)
	assert.Nil(t, res)
	var errCommon common.Error
	assert.ErrorAs(t, err, &errCommon)
	assert.Equal(t, common.ErrTypeNotFound, errCommon.Type())
}

func TestLogic_UpdateUser_StorageNilError(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	userUpdate := businesslogic.UserUpdate{
		FirstName:  "John",
		LastName:   "Doe",
		UpdateMask: []string{"first_name", "last_name"},
	}

	storageUserUpdate := storage.UserUpdate{
		FirstName: &userUpdate.FirstName,
		LastName:  &userUpdate.LastName,
	}

	ts.mockModelConverter.EXPECT().fromModelUserUpdateToStorage(gomock.Any(), userUpdate).Return(storageUserUpdate, nil)

	now := time.Now()
	ts.mockTimeProvider.EXPECT().Now().Return(now)

	storageUserUpdate.UpdatedAt = &now

	ts.mockUserStorage.EXPECT().UpdateUser(gomock.Any(), "user-id", storageUserUpdate).
		Return(nil, nil)

	res, err := ts.userManager.UpdateUser(context.Background(), "user-id", userUpdate)
	assert.Nil(t, res)
	var errCommon common.Error
	assert.ErrorAs(t, err, &errCommon)
	assert.Equal(t, common.ErrTypeInternal, errCommon.Type())
}

func TestLogic_UpdateUser_EventEmitterError(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	userUpdate := businesslogic.UserUpdate{
		FirstName:  "John",
		LastName:   "Doe",
		UpdateMask: []string{"first_name", "last_name"},
	}

	storageUserUpdate := storage.UserUpdate{
		FirstName: &userUpdate.FirstName,
		LastName:  &userUpdate.LastName,
	}

	ts.mockModelConverter.EXPECT().fromModelUserUpdateToStorage(gomock.Any(), userUpdate).Return(storageUserUpdate, nil)

	now := time.Now()
	ts.mockTimeProvider.EXPECT().Now().Return(now)

	storageUserUpdate.UpdatedAt = &now

	storageUser := &storage.User{
		ID:        "user-id",
		FirstName: *storageUserUpdate.FirstName,
		LastName:  *storageUserUpdate.LastName,
		Nickname:  "nickname",
		Email:     "email",
		Country:   "uk",
		CreatedAt: now,
		UpdatedAt: now,
	}

	ts.mockUserStorage.EXPECT().UpdateUser(gomock.Any(), "user-id", storageUserUpdate).Return(storageUser, nil)

	expectedUser := businesslogic.User{
		ID:        storageUser.ID,
		FirstName: storageUser.FirstName,
		LastName:  storageUser.LastName,
		Nickname:  storageUser.Nickname,
		Email:     storageUser.Email,
		Country:   storageUser.Country,
		CreatedAt: storageUser.CreatedAt,
		UpdatedAt: storageUser.UpdatedAt,
	}

	ts.mockModelConverter.EXPECT().fromStorageUserToModel(gomock.Any(), *storageUser).Return(expectedUser)

	userEvent := events.UserEvent{
		UserId:    storageUser.ID,
		EventType: events.EventTypeUpdated,
		EventMask: userUpdate.UpdateMask,
		EventTime: now,
	}

	ts.mockModelConverter.EXPECT().fromModelUserToEvent(gomock.Any(), expectedUser).Return(userEvent)

	ts.mockEventEmitter.EXPECT().EmitUserEvent(gomock.Any(), userEvent).
		Return(common.NewError(nil, common.ErrTypeInternal))

	res, err := ts.userManager.UpdateUser(context.Background(), "user-id", userUpdate)
	// Update should succeed even if event emission fails
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, expectedUser, *res)
}

func TestLogic_UpdateUser_Success(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	userUpdate := businesslogic.UserUpdate{
		FirstName:  "John",
		LastName:   "Doe",
		UpdateMask: []string{"first_name", "last_name"},
	}

	storageUserUpdate := storage.UserUpdate{
		FirstName: &userUpdate.FirstName,
		LastName:  &userUpdate.LastName,
	}

	ts.mockModelConverter.EXPECT().fromModelUserUpdateToStorage(gomock.Any(), userUpdate).Return(storageUserUpdate, nil)

	now := time.Now()
	ts.mockTimeProvider.EXPECT().Now().Return(now)

	storageUserUpdate.UpdatedAt = &now

	storageUser := &storage.User{
		ID:        "user-id",
		FirstName: *storageUserUpdate.FirstName,
		LastName:  *storageUserUpdate.LastName,
		Nickname:  "nickname",
		Email:     "email",
		Country:   "uk",
		CreatedAt: now,
		UpdatedAt: now,
	}

	ts.mockUserStorage.EXPECT().UpdateUser(gomock.Any(), "user-id", storageUserUpdate).Return(storageUser, nil)

	expectedUser := businesslogic.User{
		ID:        storageUser.ID,
		FirstName: storageUser.FirstName,
		LastName:  storageUser.LastName,
		Nickname:  storageUser.Nickname,
		Email:     storageUser.Email,
		Country:   storageUser.Country,
		CreatedAt: storageUser.CreatedAt,
		UpdatedAt: storageUser.UpdatedAt,
	}

	ts.mockModelConverter.EXPECT().fromStorageUserToModel(gomock.Any(), *storageUser).Return(expectedUser)

	userEvent := events.UserEvent{
		UserId:    storageUser.ID,
		EventType: events.EventTypeUpdated,
		EventMask: userUpdate.UpdateMask,
		EventTime: now,
		FirstName: storageUser.FirstName,
		LastName:  storageUser.LastName,
		Nickname:  storageUser.Nickname,
		Email:     storageUser.Email,
		Country:   storageUser.Country,
		CreatedAt: storageUser.CreatedAt,
		UpdatedAt: now,
	}

	ts.mockModelConverter.EXPECT().fromModelUserToEvent(gomock.Any(), expectedUser).Return(userEvent)

	ts.mockEventEmitter.EXPECT().EmitUserEvent(gomock.Any(), userEvent).Return(nil)

	res, err := ts.userManager.UpdateUser(context.Background(), "user-id", userUpdate)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, expectedUser, *res)
}
