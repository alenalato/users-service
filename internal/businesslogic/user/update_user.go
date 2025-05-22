package user

import (
	"context"
	"errors"
	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/events"
	"github.com/alenalato/users-service/internal/logger"
)

func (l *Logic) UpdateUser(
	ctx context.Context,
	userId string,
	userUpdate businesslogic.UserUpdate,
) (*businesslogic.User, error) {
	// Validate input
	errValidate := validate.Struct(userUpdate)
	if errValidate != nil {
		logger.Log.Errorf("validation error: %v", errValidate)

		return nil, common.NewError(errValidate, common.ErrTypeInvalidArgument)
	}

	// Prepare storage user update
	storageUserUpdate, errConv := l.converter.fromModelUserUpdateToStorage(ctx, userUpdate)
	if errConv != nil {
		return nil, errConv
	}

	// Set updated at timestamp
	now := l.time.Now().UTC()
	storageUserUpdate.UpdatedAt = &now

	// Update user in storage
	storageUser, errUpdate := l.userStorage.UpdateUser(ctx, userId, storageUserUpdate)
	if errUpdate != nil {
		return nil, errUpdate
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
	userEvent.EventType = events.EventTypeUpdated
	userEvent.EventMask = userUpdate.UpdateMask
	userEvent.EventTime = now
	errEmit := l.eventEmitter.EmitUserEvent(ctx, userEvent)
	if errEmit != nil {
		logger.Log.Warn("User updated without event emission: %v", userEvent)
	}

	return &user, nil
}
