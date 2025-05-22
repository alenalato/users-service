package user

import (
	"context"
	"github.com/alenalato/users-service/internal/events"
	"testing"
	"time"

	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestFromModelUserDetailsToStorage(t *testing.T) {
	converter := newBusinessLogicModelConverter()
	ctx := context.Background()

	model := businesslogic.UserDetails{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johnd",
		Email:     "john.doe@example.com",
		Password:  businesslogic.PasswordDetails{Hash: "hashed_password"},
		Country:   "US",
	}

	expected := storage.UserDetails{
		FirstName:    "John",
		LastName:     "Doe",
		Nickname:     "johnd",
		Email:        "john.doe@example.com",
		PasswordHash: "hashed_password",
		Country:      "US",
	}

	result := converter.fromModelUserDetailsToStorage(ctx, model)
	assert.Equal(t, expected, result)
}

func TestFromModelUserUpdateToStorage(t *testing.T) {
	converter := newBusinessLogicModelConverter()
	ctx := context.Background()

	model := businesslogic.UserUpdate{
		FirstName:  "John",
		LastName:   "Doe",
		Nickname:   "johnd",
		Email:      "john@doe.com",
		Country:    "US",
		UpdateMask: []string{"first_name", "last_name", "nickname", "email", "country"},
	}

	expected := storage.UserUpdate{
		FirstName: &model.FirstName,
		LastName:  &model.LastName,
		Nickname:  &model.Nickname,
		Email:     &model.Email,
		Country:   &model.Country,
	}

	result, err := converter.fromModelUserUpdateToStorage(ctx, model)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestFromModelUserUpdateToStorage_Partial(t *testing.T) {
	converter := newBusinessLogicModelConverter()
	ctx := context.Background()

	model := businesslogic.UserUpdate{
		FirstName:  "John",
		LastName:   "Doe",
		Country:    "US",
		UpdateMask: []string{"first_name", "country"},
	}

	expected := storage.UserUpdate{
		FirstName: &model.FirstName,
		Country:   &model.Country,
	}

	result, err := converter.fromModelUserUpdateToStorage(ctx, model)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestFromModelUserUpdateToStorage_InvalidMask(t *testing.T) {
	converter := newBusinessLogicModelConverter()
	ctx := context.Background()

	model := businesslogic.UserUpdate{
		UpdateMask: []string{"invalid_field"},
	}

	_, err := converter.fromModelUserUpdateToStorage(ctx, model)
	assert.Error(t, err)
}

func TestFromModelUserFilterToStorage(t *testing.T) {
	converter := newBusinessLogicModelConverter()
	ctx := context.Background()

	firstName := "John"
	lastName := "Doe"
	country := "US"

	model := businesslogic.UserFilter{
		FirstName: &firstName,
		LastName:  &lastName,
		Country:   &country,
	}

	expected := storage.UserFilter{
		FirstName: &firstName,
		LastName:  &lastName,
		Country:   &country,
	}

	result := converter.fromModelUserFilterToStorage(ctx, model)
	assert.Equal(t, expected, result)
}

func TestFromStorageUserToModel(t *testing.T) {
	converter := newBusinessLogicModelConverter()
	ctx := context.Background()

	storageUser := storage.User{
		ID:        "123",
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johnd",
		Email:     "john.doe@example.com",
		Country:   "US",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	expected := businesslogic.User{
		ID:        "123",
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johnd",
		Email:     "john.doe@example.com",
		Country:   "US",
		CreatedAt: storageUser.CreatedAt,
		UpdatedAt: storageUser.UpdatedAt,
	}

	result := converter.fromStorageUserToModel(ctx, storageUser)
	assert.Equal(t, expected, result)
}

func TestFromModelUserToEvent(t *testing.T) {
	converter := newBusinessLogicModelConverter()
	ctx := context.Background()

	model := businesslogic.User{
		ID:        "123",
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johnd",
		Email:     "john.doe@example.com",
		Country:   "US",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	expected := events.UserEvent{
		UserId:    "123",
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johnd",
		Email:     "john.doe@example.com",
		Country:   "US",
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}

	result := converter.fromModelUserToEvent(ctx, model)
	assert.Equal(t, expected, result)
}
