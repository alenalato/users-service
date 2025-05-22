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
	collection := testMongoStorage.Database().Collection(UserCollection)
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

	// clean up the test user
	_, errDelete := collection.DeleteOne(
		context.Background(),
		bson.D{{Key: "_id", Value: userId}},
	)
	require.NoError(t, errDelete)
}

func TestMongoDB_UpdateUser_AlreadyExistsError(t *testing.T) {
	userId := "existinguser"

	now := time.Now().UTC()

	userDetails := storage.UserDetails{
		ID:           userId,
		FirstName:    "OldFirstName",
		LastName:     "OldLastName",
		Nickname:     "nickname",
		Email:        "email",
		PasswordHash: "passwordhash",
		Country:      "country",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// create first user
	_, err := testMongoStorage.CreateUser(context.Background(), userDetails)
	require.NoError(t, err)

	userDetails2 := storage.UserDetails{
		ID:        "anotheruser",
		FirstName: "AnotherFirstName",
		LastName:  "AnotherLastName",
		Nickname:  "nickname2",
		Email:     "email2",
		Country:   "country2",
	}

	// create second user
	_, err = testMongoStorage.CreateUser(context.Background(), userDetails2)
	require.NoError(t, err)

	// attempt to update the first user with the second user's nickname
	userUpdate := storage.UserUpdate{}
	userUpdate.Nickname = &userDetails2.Nickname
	userUpdate.UpdatedAt = &now

	user, err := testMongoStorage.UpdateUser(context.Background(), userId, userUpdate)
	assert.Error(t, err)
	assert.Nil(t, user)
	var errCommon common.Error
	assert.ErrorAs(t, err, &errCommon)
	assert.Equal(t, common.ErrTypeAlreadyExists, errCommon.Type())

	// clean up the test users
	collection := testMongoStorage.Database().Collection(UserCollection)
	_, errDelete := collection.DeleteMany(
		context.Background(),
		bson.D{
			{"_id", bson.D{
				{"$in", []string{userId, userDetails2.ID}},
			}},
		},
	)
	require.NoError(t, errDelete)
}

func TestMongoDB_UpdateUser_NotFoundError(t *testing.T) {
	userId := "nonexistentuser"

	userUpdate := storage.UserUpdate{}

	updatedUser, err := testMongoStorage.UpdateUser(context.Background(), userId, userUpdate)
	assert.Error(t, err)
	assert.Nil(t, updatedUser)
	var errCommon common.Error
	assert.ErrorAs(t, err, &errCommon)
	assert.Equal(t, common.ErrTypeNotFound, errCommon.Type())
}
