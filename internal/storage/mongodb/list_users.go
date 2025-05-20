package mongodb

import (
	"context"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/logger"
	"github.com/alenalato/users-service/internal/storage"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (m *MongoDB) ListUsers(
	ctx context.Context,
	userFilter storage.UserFilter,
	pageSize int,
	pageToken string,
) ([]storage.User, string, error) {
	collection := m.database.Collection(UserCollection)

	if pageSize <= 0 || pageSize > storage.MaxPageSize {
		pageSize = storage.MaxPageSize
	}

	var skipSize int64
	if pageToken != "" {
		var tokErr error
		userFilter, skipSize, tokErr = parsePageToken(pageToken)
		if tokErr != nil {
			return nil, "", tokErr
		}
	}

	// add 1 to the requested pageSize to understand if there will be a next page, for example:
	// - 20 records on database
	// - the first request is 10 records => users[11] exists => return the nextPageToken
	// - the second request is 10 records => users[11] does not exist => don't return the nextPageToken
	// this strategy avoids a c.CountDocuments() at every request reducing database load
	limit := int64(pageSize) + 1

	opts := options.Find().SetSkip(skipSize).SetLimit(limit)

	cursor, findErr := collection.Find(ctx, userFilter, opts)
	if findErr != nil {
		logger.Log.Errorf("Error listing users: %v", findErr)

		return nil, "", common.NewError(findErr, common.ErrTypeInternal)
	}

	var users []storage.User
	if cursErr := cursor.All(ctx, &users); cursErr != nil {
		logger.Log.Errorf("Error decoding users: %v", findErr)

		return nil, "", common.NewError(findErr, common.ErrTypeInternal)
	}

	nextPageToken := ""
	if len(users) > pageSize {
		var tokGenErr error
		nextPageToken, tokGenErr = generateNextPageToken(userFilter, skipSize+int64(pageSize))
		if tokGenErr != nil {
			logger.Log.Errorf("Error generating next page token: %v", tokGenErr)

			return nil, "", common.NewError(tokGenErr, common.ErrTypeInternal)
		}

		// remove the last element from the list, as it is used only to check if there is a next page
		users = users[:pageSize]
	}

	return users, nextPageToken, nil
}
