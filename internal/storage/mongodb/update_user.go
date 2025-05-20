package mongodb

import (
	"context"
	"errors"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/logger"
	"github.com/alenalato/users-service/internal/storage"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

func (m *MongoDB) UpdateUser(ctx context.Context, userId string, userUpdate storage.UserUpdate) (*storage.User, error) {
	collection := m.database.Collection(UserCollection)

	filter := bson.D{{Key: "_id", Value: userId}}
	update := bson.D{{Key: "$set", Value: userUpdate}}
	opts := options.UpdateOne().SetUpsert(false)

	updateCtx, cancelUpdate := context.WithTimeout(ctx, 10*time.Second)
	defer cancelUpdate()

	updateRes, updateErr := collection.UpdateOne(updateCtx, filter, update, opts)
	if updateErr != nil {
		logger.Log.Errorf("Error updating user: %v", updateErr)

		return nil, common.NewError(updateErr, common.ErrTypeInternal)
	}
	if updateRes.MatchedCount == 0 {
		logger.Log.Debugf("Error updating user: %v", updateErr)

		return nil, common.NewError(errors.New("user not found"), common.ErrTypeNotFound)
	}

	findCtx, cancelFind := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFind()

	var user storage.User
	findErr := collection.FindOne(findCtx, filter).Decode(&user)
	if findErr != nil {
		logger.Log.Errorf("Error decoding updated user: %v", findErr)

		return nil, common.NewError(findErr, common.ErrTypeInternal)
	}

	return &user, nil
}
