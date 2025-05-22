package user

import (
	"context"
	"errors"
	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/events"
	"github.com/alenalato/users-service/internal/logger"
	"github.com/alenalato/users-service/internal/storage"
)

// modelConverter is an interface that defines methods for converting between
// user business logic models and other models (e.g., storage models, event models).
type modelConverter interface {
	fromModelUserDetailsToStorage(ctx context.Context, userDetails businesslogic.UserDetails) storage.UserDetails
	fromModelUserUpdateToStorage(ctx context.Context, userUpdate businesslogic.UserUpdate) (storage.UserUpdate, error)
	fromModelUserFilterToStorage(ctx context.Context, userFilter businesslogic.UserFilter) storage.UserFilter
	fromStorageUserToModel(ctx context.Context, user storage.User) businesslogic.User
	fromModelUserToEvent(ctx context.Context, user businesslogic.User) events.UserEvent
}

type businessLogicModelConverter struct {
}

var _ modelConverter = new(businessLogicModelConverter)

// fromModelUserDetailsToStorage converts a businesslogic.UserDetails to a storage.UserDetails
func (c *businessLogicModelConverter) fromModelUserDetailsToStorage(
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

// fromModelUserUpdateToStorage converts a businesslogic.UserUpdate to a storage.UserUpdate
func (c *businessLogicModelConverter) fromModelUserUpdateToStorage(
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
		case "nickname":
			storageUserUpdate.Nickname = &userUpdate.Nickname
			validUpdates = true
		case "email":
			storageUserUpdate.Email = &userUpdate.Email
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

// fromModelUserFilterToStorage converts a businesslogic.UserFilter to a storage.UserFilter
func (c *businessLogicModelConverter) fromModelUserFilterToStorage(
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

// fromStorageUserToModel converts a storage.User to a businesslogic.User
func (c *businessLogicModelConverter) fromStorageUserToModel(_ context.Context, user storage.User) businesslogic.User {
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

// fromModelUserToEvent converts a businesslogic.User to an events.UserEvent
func (c *businessLogicModelConverter) fromModelUserToEvent(
	_ context.Context,
	user businesslogic.User,
) events.UserEvent {
	return events.UserEvent{
		UserId:    user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Country:   user.Country,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// newBusinessLogicModelConverter creates a new businessLogicModelConverter instance
func newBusinessLogicModelConverter() *businessLogicModelConverter {
	return &businessLogicModelConverter{}
}
