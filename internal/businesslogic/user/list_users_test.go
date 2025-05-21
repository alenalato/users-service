package user

import (
	"context"
	"errors"
	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestLogic_ListUsers_ValidationError(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	userFilter := businesslogic.UserFilter{}
	pageSize := -1
	pageToken := ""

	users, nextPageToken, err := ts.userManager.ListUsers(context.Background(), userFilter, pageSize, pageToken)

	assert.Nil(t, users)
	assert.Empty(t, nextPageToken)
	var errCommon common.Error
	assert.ErrorAs(t, err, &errCommon)
	assert.Equal(t, common.ErrTypeInvalidArgument, errCommon.Type())
}

func TestLogic_ListUsers_StorageError(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	userFilter := businesslogic.UserFilter{}
	pageSize := 10
	pageToken := "token"

	ts.mockModelConverter.EXPECT().fromModelUserFilterToStorage(gomock.Any(), userFilter).Return(storage.UserFilter{})

	ts.mockUserStorage.EXPECT().ListUsers(gomock.Any(), storage.UserFilter{}, pageSize, pageToken).
		Return(nil, "", common.NewError(errors.New("storage error"), common.ErrTypeInternal))

	users, nextPageToken, err := ts.userManager.ListUsers(context.Background(), userFilter, pageSize, pageToken)

	assert.Nil(t, users)
	assert.Empty(t, nextPageToken)
	var errCommon common.Error
	assert.ErrorAs(t, err, &errCommon)
	assert.Equal(t, common.ErrTypeInternal, errCommon.Type())
}

func TestLogic_ListUsers_Success(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	country := "US"
	userFilter := businesslogic.UserFilter{
		Country: &country,
	}
	pageSize := 10
	pageToken := "token"

	storageUsers := []storage.User{
		{
			ID:        "user-1",
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "johndoe",
			Email:     "john@doe.com",
			Country:   "US",
		},
		{
			ID:        "user-2",
			FirstName: "Jane",
			LastName:  "Smith",
			Nickname:  "janesmith",
			Email:     "jane@smith.com",
			Country:   "US",
		},
	}
	nextPageToken := "next-token"

	storageFilter := storage.UserFilter{
		Country: &country,
	}

	ts.mockModelConverter.EXPECT().fromModelUserFilterToStorage(gomock.Any(), userFilter).Return(storageFilter)

	ts.mockUserStorage.EXPECT().ListUsers(gomock.Any(), storageFilter, pageSize, pageToken).
		Return(storageUsers, nextPageToken, nil)

	ts.mockModelConverter.EXPECT().fromStorageUserToModel(gomock.Any(), storageUsers[0]).Return(businesslogic.User{
		ID:        "user-1",
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Email:     "john@doe.com",
		Country:   "US",
	})
	ts.mockModelConverter.EXPECT().fromStorageUserToModel(gomock.Any(), storageUsers[1]).Return(businesslogic.User{
		ID:        "user-2",
		FirstName: "Jane",
		LastName:  "Smith",
		Nickname:  "janesmith",
		Email:     "jane@smith.com",
		Country:   "UK",
	})

	users, actualNextPageToken, err := ts.userManager.ListUsers(context.Background(), userFilter, pageSize, pageToken)

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)
	assert.Equal(t, nextPageToken, actualNextPageToken)
	assert.Equal(t, "user-1", users[0].ID)
	assert.Equal(t, "user-2", users[1].ID)
}
