package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/logger"
	"github.com/alenalato/users-service/internal/storage"
)

func (l *Logic) ListUsers(
	ctx context.Context,
	userFilter businesslogic.UserFilter,
	pageSize int,
	pageToken string,
) ([]businesslogic.User, string, error) {
	// validate input
	errValidate := errors.Join(
		validate.Var(pageSize, "gte=0"),
		validate.Var(pageSize, fmt.Sprintf("lte=%d", storage.MaxPageSize)),
	)
	if errValidate != nil {
		logger.Log.Errorf("validation error: %v", errValidate)

		return nil, "", common.NewError(errValidate, common.ErrTypeInvalidArgument)
	}

	storageUsers, nextPageToken, errList := l.userStorage.ListUsers(
		ctx,
		l.converter.fromModelUserFilterToStorage(ctx, userFilter),
		pageSize,
		pageToken,
	)
	if errList != nil {
		return nil, "", errList
	}

	var users []businesslogic.User
	for _, storageUser := range storageUsers {
		users = append(users, l.converter.fromStorageUserToModel(ctx, storageUser))
	}

	return users, nextPageToken, nil
}
