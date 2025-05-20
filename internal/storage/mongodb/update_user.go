package mongodb

import (
	"context"
	"errors"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/logger"
	"github.com/alenalato/users-service/internal/storage"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

func (m *MongoDB) UpdateUser(ctx context.Context, userId string, userUpdate storage.UserUpdate) (*storage.User, error) {
	collection := m.database.Collection(UserCollection)

	filter := bson.D{{Key: "_id", Value: userId}}
	update := bson.D{{Key: "$set", Value: userUpdate}}
	opts := options.FindOneAndUpdate().
		SetUpsert(false).
		SetReturnDocument(options.After)

	updateCtx, cancelUpdate := context.WithTimeout(ctx, 10*time.Second)
	defer cancelUpdate()

	var user storage.User

	updateRes := collection.FindOneAndUpdate(updateCtx, filter, update, opts)
	updateErr := updateRes.Decode(&user)
	if updateErr != nil {
		if errors.Is(updateErr, mongo.ErrNoDocuments) {
			logger.Log.Debugf("Error updating user: %v", updateErr)

			return nil, common.NewError(errors.New("user not found"), common.ErrTypeNotFound)
		} else {
			logger.Log.Errorf("Error updating user: %v", updateErr)

			return nil, common.NewError(updateErr, common.ErrTypeInternal)
		}
	}

	return &user, nil
}
