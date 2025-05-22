package mongodb

import (
	"context"
	"errors"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

func (m *MongoDB) DeleteUser(ctx context.Context, userId string) error {
	collection := m.database.Collection(UserCollection)

	deleteCtx, cancelDelete := context.WithTimeout(ctx, 10*time.Second)
	defer cancelDelete()

	filter := bson.D{{Key: "_id", Value: userId}}

	deleteRes, errDelete := collection.DeleteOne(deleteCtx, filter)
	if errDelete != nil {
		logger.Log.Errorf("Error deleting user: %v", errDelete)

		return common.NewError(errDelete, common.ErrTypeInternal)
	}
	// Check if the user was found and deleted, if not, return a not found error
	if deleteRes.DeletedCount == 0 {
		errDelete = errors.New("user not found")
		logger.Log.Errorf("Error deleting user: %v", errDelete)

		return common.NewError(errDelete, common.ErrTypeNotFound)
	}

	return nil
}
