package user

import (
	"context"
	"errors"
	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/logger"
	"github.com/alenalato/users-service/internal/storage"
)

type modelConverter interface {
	fromModelUserDetailsToStorage(ctx context.Context, userDetails businesslogic.UserDetails) storage.UserDetails
	fromModelUserUpdateToStorage(ctx context.Context, userUpdate businesslogic.UserUpdate) (storage.UserUpdate, error)
	fromModelUserFilterToStorage(ctx context.Context, userFilter businesslogic.UserFilter) storage.UserFilter
	fromStorageUserToModel(ctx context.Context, user storage.User) businesslogic.User
}

type storageModelConverter struct {
}

var _ modelConverter = new(storageModelConverter)

func (c *storageModelConverter) fromModelUserDetailsToStorage(
	_ context.Context,
	userDetails businesslogic.UserDetails,
) storage.UserDetails {
	return storage.UserDetails{
		FirstName:    userDetails.FirstName,
		LastName:     userDetails.LastName,
		Nickname:     userDetails.Nickname,
		Email:        userDetails.Email,
		PasswordHash: userDetails.Password.Hash,
		Country:      userDetails.Country,
	}
}

func (c *storageModelConverter) fromModelUserUpdateToStorage(
	_ context.Context,
	userUpdate businesslogic.UserUpdate,
) (storage.UserUpdate, error) {
	storageUserUpdate := storage.UserUpdate{}
	validUpdates := false
	for _, field := range userUpdate.UpdateMask {
		switch field {
		case "first_name":
			storageUserUpdate.FirstName = &userUpdate.FirstName
			validUpdates = true
		case "last_name":
			storageUserUpdate.LastName = &userUpdate.LastName
			validUpdates = true
		case "country":
			storageUserUpdate.Country = &userUpdate.Country
			validUpdates = true
		default:
			continue
		}
	}
	if !validUpdates {
		err := errors.New("no valid fields in update mask")
		logger.Log.Errorf("error converting user update: %s", err.Error())

		return storage.UserUpdate{}, common.NewError(err, common.ErrTypeInvalidArgument)
	}

	return storageUserUpdate, nil
}

func (c *storageModelConverter) fromModelUserFilterToStorage(
	_ context.Context,
	userFilter businesslogic.UserFilter,
) storage.UserFilter {
	storageUserFilter := storage.UserFilter{}

	if userFilter.FirstName != nil {
		storageUserFilter.FirstName = userFilter.FirstName
	}
	if userFilter.LastName != nil {
		storageUserFilter.LastName = userFilter.LastName
	}
	if userFilter.Country != nil {
		storageUserFilter.Country = userFilter.Country
	}

	return storageUserFilter
}

func (c *storageModelConverter) fromStorageUserToModel(_ context.Context, user storage.User) businesslogic.User {
	return businesslogic.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Country:   user.Country,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// newStorageModelConverter creates a new storageModelConverter
func newStorageModelConverter() *storageModelConverter {
	return &storageModelConverter{}
}
