package user

import (
	"context"
	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/internal/storage"
)

type modelConverter interface {
	fromModelUserDetailsToStorage(ctx context.Context, userDetails businesslogic.UserDetails) storage.UserDetails
	fromModelUserUpdateToStorage(ctx context.Context, userUpdate businesslogic.UserUpdate) storage.UserUpdate
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
) storage.UserUpdate {
	storageUserUpdate := storage.UserUpdate{}
	for _, field := range userUpdate.UpdateMask {
		switch field {
		case "first_name":
			storageUserUpdate.FirstName = &userUpdate.FirstName
		case "last_name":
			storageUserUpdate.LastName = &userUpdate.LastName
		case "country":
			storageUserUpdate.Country = &userUpdate.Country
		default:
			continue
		}
	}

	return storageUserUpdate
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
