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
		var errTok error
		userFilter, skipSize, errTok = parsePageToken(pageToken)
		if errTok != nil {
			return nil, "", errTok
		}
	}

	// add 1 to the requested pageSize to test if there will be a next page
	limit := int64(pageSize) + 1

	opts := options.Find().SetSkip(skipSize).SetLimit(limit)

	cursor, errFind := collection.Find(ctx, userFilter, opts)
	if errFind != nil {
		logger.Log.Errorf("Error listing users: %v", errFind)

		return nil, "", common.NewError(errFind, common.ErrTypeInternal)
	}

	var users []storage.User
	if errCurs := cursor.All(ctx, &users); errCurs != nil {
		logger.Log.Errorf("Error decoding users: %v", errFind)

		return nil, "", common.NewError(errFind, common.ErrTypeInternal)
	}

	nextPageToken := ""
	// there is the extra element in the list, generate the next page token
	if len(users) > pageSize {
		var errTokGen error
		nextPageToken, errTokGen = generateNextPageToken(userFilter, skipSize+int64(pageSize))
		if errTokGen != nil {
			logger.Log.Errorf("Error generating next page token: %v", errTokGen)

			return nil, "", common.NewError(errTokGen, common.ErrTypeInternal)
		}

		// remove the last element from the list
		users = users[:pageSize]
	}

	return users, nextPageToken, nil
}
