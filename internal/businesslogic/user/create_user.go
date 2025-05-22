package user

import (
	"context"
	"errors"
	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/events"
	"github.com/alenalato/users-service/internal/logger"
	"github.com/google/uuid"
)

func (l *Logic) CreateUser(ctx context.Context, userDetails businesslogic.UserDetails) (*businesslogic.User, error) {
	// Validate input
	errValidate := validate.Struct(userDetails)
	if errValidate != nil {
		logger.Log.Errorf("validation error: %v", errValidate)

		return nil, common.NewError(errValidate, common.ErrTypeInvalidArgument)
	}

	// Hash password using password manager
	passwordErr := l.passwordManager.GeneratePasswordHash(ctx, &userDetails.Password)
	if passwordErr != nil {
		return nil, passwordErr
	}

	// Prepare storage user details
	storageUserDetails := l.converter.fromModelUserDetailsToStorage(ctx, userDetails)

	// Generate ID and timestamps
	storageUserDetails.ID = uuid.New().String()
	now := l.time.Now().UTC()
	storageUserDetails.CreatedAt = now
	storageUserDetails.UpdatedAt = now

	// Create user in storage
	storageUser, errCreate := l.userStorage.CreateUser(ctx, storageUserDetails)
	if errCreate != nil {
		return nil, errCreate
	}
	if storageUser == nil {
		err := errors.New("unexpected nil storage user")
		logger.Log.Error(err)

		return nil, common.NewError(err, common.ErrTypeInternal)
	}

	// Convert storage user to model user
	user := l.converter.fromStorageUserToModel(ctx, *storageUser)

	// Emit user event
	userEvent := l.converter.fromModelUserToEvent(ctx, user)
	userEvent.EventType = events.EventTypeCreated
	userEvent.EventTime = now
	errEmit := l.eventEmitter.EmitUserEvent(ctx, userEvent)
	if errEmit != nil {
		logger.Log.Warn("User created without event emission: %v", userEvent)
	}

	return &user, nil
}
