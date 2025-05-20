package user

import (
	"context"
	"errors"
	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/logger"
	"github.com/google/uuid"
)

func (l *Logic) CreateUser(ctx context.Context, userDetails businesslogic.UserDetails) (*businesslogic.User, error) {
	// validate input
	validateErr := validate.Struct(userDetails)
	if validateErr != nil {
		logger.Log.Errorf("validation error: %v", validateErr)

		return nil, common.NewError(validateErr, common.ErrTypeInvalidArgument)
	}

	// hash password
	passwordErr := l.passwordManager.GeneratePasswordHash(ctx, &userDetails.Password)
	if passwordErr != nil {
		return nil, passwordErr
	}

	// prepare storage user details
	storageUserDetails := l.converter.fromModelUserDetailsToStorage(ctx, userDetails)

	// generate ID and timestamps
	storageUserDetails.ID = uuid.New().String()
	now := l.time.Now()
	storageUserDetails.CreatedAt = now
	storageUserDetails.UpdatedAt = now

	// create user in storage
	storageUser, createErr := l.userStorage.CreateUser(ctx, storageUserDetails)
	if createErr != nil {
		return nil, createErr
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
