package mongodb

import (
	"context"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"testing"
)

func TestMongoDB_DeleteUser_NotFound(t *testing.T) {
	userId := "notpresent"
	err := testMongoStorage.DeleteUser(context.Background(), userId)

	assert.Error(t, err)
	var errCommon common.Error
	assert.ErrorAs(t, err, &errCommon)
	assert.Equal(t, common.ErrTypeNotFound, errCommon.Type())
}

func TestMongoDB_DeleteUser_Success(t *testing.T) {
	userId := "present"

	// create a user to ensure it exists
	collection := testMongoStorage.database.Collection(UserCollection)

	_, errIns := collection.InsertOne(
		context.Background(),
		bson.D{
			{Key: "_id", Value: userId},
		},
	)
	require.NoError(t, errIns)

	err := testMongoStorage.DeleteUser(context.Background(), userId)
	assert.NoError(t, err)

	// verify that the user is actually deleted
	user := &storage.User{}
	errFind := collection.FindOne(
		context.Background(),
		bson.D{
			{Key: "_id", Value: userId},
		},
	).Decode(user)
	assert.Error(t, errFind)
	assert.ErrorIs(t, errFind, mongo.ErrNoDocuments)
}
