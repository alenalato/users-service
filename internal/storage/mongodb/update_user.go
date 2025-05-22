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
		SetUpsert(false).                // Do not create a new document if the filter does not match
		SetReturnDocument(options.After) // Return the updated document

	updateCtx, cancelUpdate := context.WithTimeout(ctx, 10*time.Second)
	defer cancelUpdate()

	var user storage.User

	// Update the user in the database and return the updated document in a single operation
	updateRes := collection.FindOneAndUpdate(updateCtx, filter, update, opts)
	updateErr := updateRes.Decode(&user)
	if updateErr != nil {
		if errors.Is(updateErr, mongo.ErrNoDocuments) { // Check if the error is due to the user not being found
			logger.Log.Debugf("Error updating user: %v", updateErr)

			return nil, common.NewError(errors.New("user not found"), common.ErrTypeNotFound)
		} else if mongo.IsDuplicateKeyError(updateErr) { // Check for duplicate key error
			logger.Log.Debugf("Error updating user: %v", updateErr)

			return nil, common.NewError(
				errors.New("another user with same nickname or email already exists"),
				common.ErrTypeAlreadyExists,
			)
		} else {
			logger.Log.Errorf("Error updating user: %v", updateErr)

			return nil, common.NewError(updateErr, common.ErrTypeInternal)
		}
	}

	return &user, nil
}
