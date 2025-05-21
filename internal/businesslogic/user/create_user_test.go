package user

import (
	"context"
	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/events"
	"github.com/alenalato/users-service/internal/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

type storageUserDetailsMatcher struct {
	expected storage.UserDetails
}

func (m storageUserDetailsMatcher) Matches(x interface{}) bool {
	actual, ok := x.(storage.UserDetails)
	if !ok {
		return false
	}

	return cmp.Equal(m.expected, actual,
		cmpopts.IgnoreFields(storage.UserDetails{}, "ID")) &&
		actual.ID != ""
}

func (m storageUserDetailsMatcher) String() string {
	return "matches storage.UserDetails"
}

func TestLogic_CreateUser_ValidationError(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	userDetails := businesslogic.UserDetails{}

	res, err := ts.userManager.CreateUser(context.Background(), userDetails)
	assert.Nil(t, res)
	var errCommon common.Error
	assert.ErrorAs(t, err, &errCommon)
	assert.Equal(t, common.ErrTypeInvalidArgument, errCommon.Type())
}

func TestLogic_CreateUser_PasswordHashError(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	userDetails := businesslogic.UserDetails{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Email:     "john@doe.it",
		Password: businesslogic.PasswordDetails{
			Text: "password",
		},
		Country: "uk",
	}

	ts.mockPasswordManager.EXPECT().GeneratePasswordHash(gomock.Any(), &userDetails.Password).
		Do(func(ctx context.Context, password *businesslogic.PasswordDetails) {
			password.Hash = "hashed_password"
		}).
		Return(common.NewError(nil, common.ErrTypeInternal))

	res, err := ts.userManager.CreateUser(context.Background(), userDetails)
	assert.Nil(t, res)
	var errCommon common.Error
	assert.ErrorAs(t, err, &errCommon)
	assert.Equal(t, common.ErrTypeInternal, errCommon.Type())
}

func TestLogic_CreateUser_StorageError(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	userDetails := businesslogic.UserDetails{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Email:     "john@doe.it",
		Password: businesslogic.PasswordDetails{
			Text: "password",
		},
		Country: "uk",
	}

	ts.mockPasswordManager.EXPECT().GeneratePasswordHash(gomock.Any(), &userDetails.Password).
		Do(func(ctx context.Context, password *businesslogic.PasswordDetails) {
			password.Hash = "hashed_password"
		}).
		Return(nil)

	storageUserDetails := storage.UserDetails{
		FirstName:    userDetails.FirstName,
		LastName:     userDetails.LastName,
		Nickname:     userDetails.Nickname,
		Email:        userDetails.Email,
		PasswordHash: "hashed_password",
		Country:      userDetails.Country,
	}

	userDetailsWithPasswordHash := userDetails
	userDetailsWithPasswordHash.Password.Hash = "hashed_password"

	ts.mockModelConverter.EXPECT().fromModelUserDetailsToStorage(gomock.Any(), userDetailsWithPasswordHash).Return(storageUserDetails)

	now := time.Now().UTC()
	ts.mockTimeProvider.EXPECT().Now().Return(now)

	storageUserDetailsWitTimestamps := storageUserDetails
	storageUserDetailsWitTimestamps.CreatedAt = now
	storageUserDetailsWitTimestamps.UpdatedAt = now

	ts.mockUserStorage.EXPECT().CreateUser(gomock.Any(), storageUserDetailsMatcher{storageUserDetailsWitTimestamps}).
		Return(nil, common.NewError(nil, common.ErrTypeAlreadyExists))

	res, err := ts.userManager.CreateUser(context.Background(), userDetails)

	assert.Nil(t, res)
	assert.Error(t, err)
	var errCommon common.Error
	assert.ErrorAs(t, err, &errCommon)
	assert.Equal(t, common.ErrTypeAlreadyExists, errCommon.Type())
}

func TestLogic_CreateUser_StorageNilUserError(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	userDetails := businesslogic.UserDetails{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Email:     "john@doe.it",
		Password: businesslogic.PasswordDetails{
			Text: "password",
		},
		Country: "uk",
	}

	ts.mockPasswordManager.EXPECT().GeneratePasswordHash(gomock.Any(), &userDetails.Password).
		Do(func(ctx context.Context, password *businesslogic.PasswordDetails) {
			password.Hash = "hashed_password"
		}).
		Return(nil)

	storageUserDetails := storage.UserDetails{
		FirstName:    userDetails.FirstName,
		LastName:     userDetails.LastName,
		Nickname:     userDetails.Nickname,
		Email:        userDetails.Email,
		PasswordHash: userDetails.Password.Hash,
		Country:      userDetails.Country,
	}

	userDetailsWithPasswordHash := userDetails
	userDetailsWithPasswordHash.Password.Hash = "hashed_password"

	ts.mockModelConverter.EXPECT().fromModelUserDetailsToStorage(gomock.Any(), userDetailsWithPasswordHash).Return(storageUserDetails)

	now := time.Now().UTC()
	ts.mockTimeProvider.EXPECT().Now().Return(now)

	storageUserDetailsWitTimestamps := storageUserDetails
	storageUserDetailsWitTimestamps.CreatedAt = now
	storageUserDetailsWitTimestamps.UpdatedAt = now

	ts.mockUserStorage.EXPECT().CreateUser(gomock.Any(), storageUserDetailsMatcher{storageUserDetailsWitTimestamps}).
		Return(nil, nil)

	res, err := ts.userManager.CreateUser(context.Background(), userDetails)

	assert.Nil(t, res)
	assert.Error(t, err)
	var errCommon common.Error
	assert.ErrorAs(t, err, &errCommon)
	assert.Equal(t, common.ErrTypeInternal, errCommon.Type())
}

func TestLogic_CreateUser_EventEmitterError(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	userDetails := businesslogic.UserDetails{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Email:     "john@doe.it",
		Password: businesslogic.PasswordDetails{
			Text: "password",
		},
		Country: "uk",
	}

	ts.mockPasswordManager.EXPECT().GeneratePasswordHash(gomock.Any(), &userDetails.Password).
		Do(func(ctx context.Context, password *businesslogic.PasswordDetails) {
			password.Hash = "hashed_password"
		}).
		Return(nil)

	storageUserDetails := storage.UserDetails{
		FirstName:    userDetails.FirstName,
		LastName:     userDetails.LastName,
		Nickname:     userDetails.Nickname,
		Email:        userDetails.Email,
		PasswordHash: userDetails.Password.Hash,
		Country:      userDetails.Country,
	}

	userDetailsWithPasswordHash := userDetails
	userDetailsWithPasswordHash.Password.Hash = "hashed_password"

	ts.mockModelConverter.EXPECT().fromModelUserDetailsToStorage(gomock.Any(), userDetailsWithPasswordHash).Return(storageUserDetails)

	now := time.Now().UTC()
	ts.mockTimeProvider.EXPECT().Now().Return(now)

	storageUser := &storage.User{
		FirstName: storageUserDetails.FirstName,
		LastName:  storageUserDetails.LastName,
		Nickname:  storageUserDetails.Nickname,
		Email:     storageUserDetails.Email,
		Country:   storageUserDetails.Country,
		CreatedAt: now,
		UpdatedAt: now,
	}

	storageUserDetailsWitTimestamps := storageUserDetails
	storageUserDetailsWitTimestamps.CreatedAt = now
	storageUserDetailsWitTimestamps.UpdatedAt = now

	var generatedUserID string
	ts.mockUserStorage.EXPECT().CreateUser(gomock.Any(), storageUserDetailsMatcher{storageUserDetailsWitTimestamps}).
		Do(func(ctx context.Context, userDetails storage.UserDetails) {
			generatedUserID = userDetails.ID
		}).
		Return(storageUser, nil)

	storageUser.ID = generatedUserID

	expectedUser := businesslogic.User{
		ID:        generatedUserID,
		FirstName: userDetails.FirstName,
		LastName:  userDetails.LastName,
		Nickname:  userDetails.Nickname,
		Email:     userDetails.Email,
		Country:   userDetails.Country,
		CreatedAt: now,
		UpdatedAt: now,
	}

	ts.mockModelConverter.EXPECT().fromStorageUserToModel(gomock.Any(), *storageUser).Return(expectedUser)

	userEvent := events.UserEvent{
		UserId:    generatedUserID,
		EventType: events.EventTypeCreated,
		EventTime: now,
		FirstName: userDetails.FirstName,
		LastName:  userDetails.LastName,
		Nickname:  userDetails.Nickname,
		Email:     userDetails.Email,
		Country:   userDetails.Country,
		CreatedAt: now,
		UpdatedAt: now,
	}

	ts.mockModelConverter.EXPECT().fromModelUserToEvent(gomock.Any(), expectedUser).Return(userEvent)

	ts.mockEventEmitter.EXPECT().EmitUserEvent(gomock.Any(), userEvent).
		Return(common.NewError(nil, common.ErrTypeInternal))

	res, err := ts.userManager.CreateUser(context.Background(), userDetails)
	assert.NoError(t, err) // Creation should succeed even if event emission fails
	assert.NotNil(t, res)
	assert.Equal(t, expectedUser, *res)
}

func TestLogic_CreateUser_Success(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	userDetails := businesslogic.UserDetails{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Email:     "john@doe.it",
		Password: businesslogic.PasswordDetails{
			Text: "password",
		},
		Country: "uk",
	}

	ts.mockPasswordManager.EXPECT().GeneratePasswordHash(gomock.Any(), &userDetails.Password).
		Do(func(ctx context.Context, password *businesslogic.PasswordDetails) {
			password.Hash = "hashed_password"
		}).
		Return(nil)

	storageUserDetails := storage.UserDetails{
		FirstName:    userDetails.FirstName,
		LastName:     userDetails.LastName,
		Nickname:     userDetails.Nickname,
		Email:        userDetails.Email,
		PasswordHash: userDetails.Password.Hash,
		Country:      userDetails.Country,
	}

	userDetailsWithPasswordHash := userDetails
	userDetailsWithPasswordHash.Password.Hash = "hashed_password"

	ts.mockModelConverter.EXPECT().fromModelUserDetailsToStorage(gomock.Any(), userDetailsWithPasswordHash).Return(storageUserDetails)

	now := time.Now().UTC()
	ts.mockTimeProvider.EXPECT().Now().Return(now)

	storageUser := &storage.User{
		FirstName: storageUserDetails.FirstName,
		LastName:  storageUserDetails.LastName,
		Nickname:  storageUserDetails.Nickname,
		Email:     storageUserDetails.Email,
		Country:   storageUserDetails.Country,
		CreatedAt: now,
		UpdatedAt: now,
	}

	storageUserDetailsWitTimestamps := storageUserDetails
	storageUserDetailsWitTimestamps.CreatedAt = now
	storageUserDetailsWitTimestamps.UpdatedAt = now

	var generatedUserID string
	ts.mockUserStorage.EXPECT().CreateUser(gomock.Any(), storageUserDetailsMatcher{storageUserDetailsWitTimestamps}).
		Do(func(ctx context.Context, userDetails storage.UserDetails) {
			generatedUserID = userDetails.ID
		}).
		Return(storageUser, nil)

	storageUser.ID = generatedUserID

	expectedUser := businesslogic.User{
		ID:        generatedUserID,
		FirstName: userDetails.FirstName,
		LastName:  userDetails.LastName,
		Nickname:  userDetails.Nickname,
		Email:     userDetails.Email,
		Country:   userDetails.Country,
		CreatedAt: now,
		UpdatedAt: now,
	}

	ts.mockModelConverter.EXPECT().fromStorageUserToModel(gomock.Any(), *storageUser).Return(expectedUser)

	userEvent := events.UserEvent{
		UserId:    generatedUserID,
		EventType: events.EventTypeCreated,
		EventTime: now,
		FirstName: userDetails.FirstName,
		LastName:  userDetails.LastName,
		Nickname:  userDetails.Nickname,
		Email:     userDetails.Email,
		Country:   userDetails.Country,
		CreatedAt: now,
		UpdatedAt: now,
	}

	ts.mockModelConverter.EXPECT().fromModelUserToEvent(gomock.Any(), expectedUser).Return(userEvent)

	ts.mockEventEmitter.EXPECT().EmitUserEvent(gomock.Any(), userEvent).Return(nil)

	res, err := ts.userManager.CreateUser(context.Background(), userDetails)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, expectedUser, *res)
}
