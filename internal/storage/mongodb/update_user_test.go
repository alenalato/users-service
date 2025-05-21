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

func TestMongoDB_UpdateUser_Success(t *testing.T) {
	userId := "existinguser"

	createdAt := time.Now().UTC()

	userDetails := storage.UserDetails{
		ID:           userId,
		FirstName:    "OldFirstName",
		LastName:     "OldLastName",
		Nickname:     "nickname",
		Email:        "email",
		PasswordHash: "passwordhash",
		Country:      "country",
		CreatedAt:    createdAt,
		UpdatedAt:    createdAt,
	}

	// create a user to ensure it exists
	collection := testMongoStorage.database.Collection(UserCollection)
	_, errIns := collection.InsertOne(
		context.Background(),
		userDetails,
	)
	require.NoError(t, errIns)

	// update user details
	firstName := "NewFirstName"
	lastName := "NewLastName"
	updatedAt := createdAt.Add(time.Hour)
	userUpdate := storage.UserUpdate{
		FirstName: &firstName,
		LastName:  &lastName,
		UpdatedAt: &updatedAt,
	}

	updatedUser, err := testMongoStorage.UpdateUser(context.Background(), userId, userUpdate)
	require.NoError(t, err)

	assert.NotNil(t, updatedUser)
	assert.Equal(t, userId, updatedUser.ID)
	assert.Equal(t, *userUpdate.FirstName, updatedUser.FirstName)
	assert.Equal(t, *userUpdate.LastName, updatedUser.LastName)
	assert.Equal(t, userDetails.Nickname, updatedUser.Nickname)
	assert.Equal(t, userDetails.Email, updatedUser.Email)
	assert.Equal(t, userDetails.Country, updatedUser.Country)
	assert.Equal(t, userDetails.CreatedAt.UnixMilli(), updatedUser.CreatedAt.UnixMilli())
	assert.Equal(t, userUpdate.UpdatedAt.UnixMilli(), updatedUser.UpdatedAt.UnixMilli())

	// verify the user is updated in the database
	var foundUser storage.User
	errFind := collection.FindOne(
		context.Background(),
		bson.D{{Key: "_id", Value: userId}},
	).Decode(&foundUser)
	assert.NoError(t, errFind)
	assert.Equal(t, *userUpdate.FirstName, foundUser.FirstName)
	assert.Equal(t, *userUpdate.LastName, foundUser.LastName)
	assert.Equal(t, userDetails.Nickname, foundUser.Nickname)
	assert.Equal(t, userDetails.Email, foundUser.Email)
	assert.Equal(t, userDetails.Country, foundUser.Country)
	assert.Equal(t, userDetails.CreatedAt.UnixMilli(), foundUser.CreatedAt.UnixMilli())
	assert.Equal(t, userUpdate.UpdatedAt.UnixMilli(), foundUser.UpdatedAt.UnixMilli())
}

func TestMongoDB_UpdateUser_NotFound(t *testing.T) {
	userId := "nonexistentuser"

	userUpdate := storage.UserUpdate{}

	updatedUser, err := testMongoStorage.UpdateUser(context.Background(), userId, userUpdate)
	assert.Error(t, err)
	assert.Nil(t, updatedUser)
	var errCommon common.Error
	assert.ErrorAs(t, err, &errCommon)
	assert.Equal(t, common.ErrTypeNotFound, errCommon.Type())
}
