package mongodb

import (
	"context"
	"errors"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/logger"
	"github.com/alenalato/users-service/internal/storage"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

func (m *MongoDB) CreateUser(ctx context.Context, userDetails storage.UserDetails) (*storage.User, error) {
	collection := m.database.Collection(UserCollection)

	insertCtx, cancelInsert := context.WithTimeout(ctx, 10*time.Second)
	defer cancelInsert()

	insertRes, insertErr := collection.InsertOne(insertCtx, userDetails)
	if insertErr != nil {
		if mongo.IsDuplicateKeyError(insertErr) {
			logger.Log.Debugf("Error creating user: %v", insertErr)
			insertErr = common.NewError(
				errors.New("another user with same nickname or email already exists"),
				common.ErrTypeAlreadyExists,
			)
		} else {
			logger.Log.Errorf("Error creating user: %v", insertErr)
			insertErr = common.NewError(insertErr, common.ErrTypeInternal)
		}

		return nil, insertErr
	}

	findCtx, cancelFind := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFind()

	filter := bson.D{{Key: "_id", Value: insertRes.InsertedID}}

	var user storage.User
	findErr := collection.FindOne(findCtx, filter).Decode(&user)
	if findErr != nil {
		logger.Log.Errorf("Error decoding created user: %v", findErr)

		return nil, common.NewError(findErr, common.ErrTypeInternal)
	}

	return &user, nil
}
