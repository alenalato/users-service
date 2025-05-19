package mongodb

import (
	"context"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/logger"
	"github.com/alenalato/users-service/internal/storage"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

func (m *MongoDB) UpdateUser(ctx context.Context, userId string, userUpdate storage.UserUpdate) (*storage.User, error) {
	updateCtx, cancelUpdate := context.WithTimeout(ctx, 10*time.Second)
	defer cancelUpdate()

	collection := m.database.Collection(UserCollection)

	filter := bson.D{{Key: "_id", Value: userId}}
	update := bson.D{{Key: "$set", Value: userUpdate}}
	opts := options.UpdateOne().SetUpsert(false)

	updateRes, updateErr := collection.UpdateOne(updateCtx, filter, update, opts)
	if updateErr != nil {
		updateErr = common.NewError(updateErr, common.ErrTypeInternal)
		logger.Log.Errorf("Error updating user: %s", updateErr.Error())

		return nil, updateErr
	}
	if updateRes.MatchedCount == 0 {
		updateErr = common.NewError(updateErr, common.ErrTypeNotFound)
		logger.Log.Debugf("Not found while updating user: %s", updateErr.Error())

		return nil, updateErr
	}

	findCtx, cancelFind := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFind()

	var user storage.User
	findErr := collection.FindOne(findCtx, filter).Decode(&user)
	if findErr != nil {
		logger.Log.Errorf("Error decoding updated user: %s", findErr.Error())

		return nil, common.NewError(findErr, common.ErrTypeInternal)
	}

	return &user, nil
}
