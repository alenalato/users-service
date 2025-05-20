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

	deleteRes, deleteErr := collection.DeleteOne(deleteCtx, filter)
	if deleteErr != nil {
		logger.Log.Errorf("Error deleting user: %v", deleteErr)

		return common.NewError(deleteErr, common.ErrTypeInternal)
	}
	if deleteRes.DeletedCount == 0 {
		deleteErr = errors.New("user not found")
		logger.Log.Errorf("Error deleting user: %v", deleteErr)

		return common.NewError(deleteErr, common.ErrTypeNotFound)
	}

	return nil
}
