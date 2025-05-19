package user

import (
	"context"
	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/internal/storage"
)

type modelConverter interface {
	fromModelUserDetailsToStorage(ctx context.Context, userDetails businesslogic.UserDetails) storage.UserDetails
	fromStorageUserToModel(ctx context.Context, user storage.User) businesslogic.User
}

type storageModelConverter struct {
}

var _ modelConverter = new(storageModelConverter)

func (c *storageModelConverter) fromModelUserDetailsToStorage(_ context.Context, userDetails businesslogic.UserDetails) storage.UserDetails {
	return storage.UserDetails{
		FirstName:    userDetails.FirstName,
		LastName:     userDetails.LastName,
		Nickname:     userDetails.Nickname,
		Email:        userDetails.Email,
		PasswordHash: userDetails.Password.Hash,
		Country:      userDetails.Country,
	}
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
