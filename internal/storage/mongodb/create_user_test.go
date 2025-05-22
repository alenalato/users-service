package mongodb

import (
	"context"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"testing"
	"time"
)

func TestMongoDB_CreateUser_Success(t *testing.T) {
	userDetails := storage.UserDetails{
		ID:           "testuser",
		FirstName:    "Test",
		LastName:     "User",
		Nickname:     "testuser",
		Email:        "testuser@example.com",
		PasswordHash: "testpassword",
		Country:      "UK",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	createdUser, err := testMongoStorage.CreateUser(context.Background(), userDetails)
	require.NoError(t, err)

	assert.NotNil(t, createdUser)
	assert.Equal(t, userDetails.ID, createdUser.ID)
	assert.Equal(t, userDetails.FirstName, createdUser.FirstName)
	assert.Equal(t, userDetails.LastName, createdUser.LastName)
	assert.Equal(t, userDetails.Nickname, createdUser.Nickname)
	assert.Equal(t, userDetails.Email, createdUser.Email)
	assert.Equal(t, userDetails.Country, createdUser.Country)
	assert.Equal(t, userDetails.CreatedAt.UnixMilli(), createdUser.CreatedAt.UnixMilli())
	assert.Equal(t, userDetails.UpdatedAt.UnixMilli(), createdUser.UpdatedAt.UnixMilli())

	// verify the user exists in the database
	collection := testMongoStorage.Database().Collection(UserCollection)
	var foundUser storage.User
	errFind := collection.FindOne(
		context.Background(),
		bson.D{{Key: "_id", Value: createdUser.ID}},
	).Decode(&foundUser)
	assert.NoError(t, errFind)
	assert.Equal(t, createdUser.ID, foundUser.ID)
	assert.Equal(t, createdUser.FirstName, foundUser.FirstName)
	assert.Equal(t, createdUser.LastName, foundUser.LastName)
	assert.Equal(t, createdUser.Nickname, foundUser.Nickname)
	assert.Equal(t, createdUser.Email, foundUser.Email)
	assert.Equal(t, createdUser.Country, foundUser.Country)
	assert.Equal(t, createdUser.CreatedAt.UnixMilli(), foundUser.CreatedAt.UnixMilli())
	assert.Equal(t, createdUser.UpdatedAt.UnixMilli(), foundUser.UpdatedAt.UnixMilli())

	// clean up the test user
	_, errDelete := collection.DeleteOne(
		context.Background(),
		bson.D{{Key: "_id", Value: createdUser.ID}},
	)
	require.NoError(t, errDelete)
}

func TestMongoDB_CreateUser_AlreadyExistsError(t *testing.T) {
	userDetails := storage.UserDetails{
		ID: "duplicateuser",
	}

	_, err := testMongoStorage.CreateUser(context.Background(), userDetails)
	assert.NoError(t, err)

	// attempt to create the same user again
	_, errDup := testMongoStorage.CreateUser(context.Background(), userDetails)
	assert.Error(t, errDup)
	var errCommon common.Error
	assert.ErrorAs(t, errDup, &errCommon)
	assert.Equal(t, common.ErrTypeAlreadyExists, errCommon.Type())

	// clean up the test user
	collection := testMongoStorage.Database().Collection(UserCollection)
	_, errDelete := collection.DeleteOne(
		context.Background(),
		bson.D{{Key: "_id", Value: userDetails.ID}},
	)
	require.NoError(t, errDelete)
}
