package user

import (
	"context"
	"errors"
	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/logger"
	"time"
)

func (l *Logic) UpdateUser(
	ctx context.Context,
	userId string,
	userUpdate businesslogic.UserUpdate,
) (*businesslogic.User, error) {
	// validate input
	validateErr := validate.Struct(userUpdate)
	if validateErr != nil {
		logger.Log.Errorf("validation error: %v", validateErr)

		return nil, common.NewError(validateErr, common.ErrTypeInvalidArgument)
	}

	// prepare storage user update
	storageUserUpdate := l.converter.fromModelUserUpdateToStorage(ctx, userUpdate)
	if storageUserUpdate.FirstName == nil && storageUserUpdate.LastName == nil && storageUserUpdate.Country == nil {
		err := errors.New("no valid fields in update mask")
		logger.Log.Error(err)

		return nil, common.NewError(err, common.ErrTypeInvalidArgument)
	}

	// set updated at timestamp
	now := time.Now()
	storageUserUpdate.UpdatedAt = &now

	storageUser, updateErr := l.userStorage.UpdateUser(ctx, userId, storageUserUpdate)
	if updateErr != nil {
		return nil, updateErr
	}
	if storageUser == nil {
		err := errors.New("unexpected nil storage user")
		logger.Log.Error(err)

		return nil, common.NewError(err, common.ErrTypeInternal)
	}

	// convert storage user to model user
	user := l.converter.fromStorageUserToModel(ctx, *storageUser)

	return &user, nil
}
