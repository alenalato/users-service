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

	// Coerce the page size to a valid value
	if pageSize <= 0 || pageSize > storage.MaxPageSize {
		pageSize = storage.MaxPageSize
	}

	var skipSize int64
	if pageToken != "" {
		// Decode the page token to get cursor skip size and filter
		var errTok error
		userFilter, skipSize, errTok = parsePageToken(pageToken)
		if errTok != nil {
			return nil, "", errTok
		}
	}

	// Add 1 to the requested pageSize to tell if there is a next page
	limit := int64(pageSize) + 1

	opts := options.Find().
		SetSkip(skipSize).
		SetLimit(limit).
		SetSort(map[string]interface{}{"created_at": 1}) // Default fixed sort by created_at in ascending order

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
	// There is the extra element at the end of the result, generate the next page token
	if len(users) > pageSize {
		var errTokGen error
		nextPageToken, errTokGen = generateNextPageToken(
			userFilter,
			skipSize+int64(pageSize), // Move the skip size to the next page
		)
		if errTokGen != nil {
			logger.Log.Errorf("Error generating next page token: %v", errTokGen)

			return nil, "", common.NewError(errTokGen, common.ErrTypeInternal)
		}

		// Remove the last element from the list
		users = users[:pageSize]
	}

	return users, nextPageToken, nil
}
