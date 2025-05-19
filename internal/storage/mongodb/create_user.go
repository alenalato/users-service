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
	insertCtx, cancelInsert := context.WithTimeout(ctx, 10*time.Second)
	defer cancelInsert()

	collection := m.database.Collection(UserCollection)

	insertRes, insertErr := collection.InsertOne(insertCtx, userDetails)
	if insertErr != nil {
		if mongo.IsDuplicateKeyError(insertErr) {
			insertErr = common.NewError(
				errors.New("another user with same nickname or email already exists"),
				common.ErrTypeAlreadyExists,
			)
		} else {
			insertErr = common.NewError(insertErr, common.ErrTypeInternal)
		}
		logger.Log.Errorf("Error creating user: %s", insertErr.Error())

		return nil, insertErr
	}

	findCtx, cancelFind := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFind()

	filter := bson.D{{Key: "_id", Value: insertRes.InsertedID}}

	var user storage.User
	findErr := collection.FindOne(findCtx, filter).Decode(&user)
	if findErr != nil {
		logger.Log.Errorf("Error decoding created user: %s", findErr.Error())

		return nil, common.NewError(findErr, common.ErrTypeInternal)
	}

	return &user, nil
}
