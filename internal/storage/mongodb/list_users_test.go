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

func testTimeForStorage(t time.Time) time.Time {
	return time.UnixMilli(t.UnixMilli()).UTC()
}

func TestMongoDB_ListUsers_Success(t *testing.T) {
	collection := testMongoStorage.Database().Collection(UserCollection)
	users := []interface{}{
		storage.User{
			ID:        "user1",
			FirstName: "Alice",
			LastName:  "Smith",
			Email:     "alice@example.com",
			Nickname:  "alice",
			Country:   "USA",
			CreatedAt: testTimeForStorage(time.Now().Add(-4 * time.Hour)),
			UpdatedAt: testTimeForStorage(time.Now().Add(-3 * time.Hour)),
		},
		storage.User{
			ID:        "user2",
			FirstName: "Bob",
			LastName:  "Brown",
			Email:     "bob@example.com",
			Nickname:  "bobby",
			Country:   "Canada",
			CreatedAt: testTimeForStorage(time.Now().Add(-3 * time.Hour)),
			UpdatedAt: testTimeForStorage(time.Now().Add(-2 * time.Hour)),
		},
		storage.User{
			ID:        "user3",
			FirstName: "Charlie",
			LastName:  "Davis",
			Email:     "charlie@example.com",
			Nickname:  "charlie",
			Country:   "UK",
			CreatedAt: testTimeForStorage(time.Now().Add(-2 * time.Hour)),
			UpdatedAt: testTimeForStorage(time.Now().Add(-1 * time.Hour)),
		},
	}
	_, err := collection.InsertMany(context.Background(), users)
	require.NoError(t, err)

	// Test ListUsers
	userFilter := storage.UserFilter{}
	pageSize := 2
	pageToken := ""

	result, nextPageToken, err := testMongoStorage.ListUsers(context.Background(), userFilter, pageSize, pageToken)
	require.NoError(t, err)
	assert.Len(t, result, pageSize)
	assert.NotEmpty(t, nextPageToken)
	assert.Equal(t, users[0], result[0])
	assert.Equal(t, users[1], result[1])

	// Test ListUsers with next page token
	pageToken = nextPageToken
	result, nextPageToken, err = testMongoStorage.ListUsers(context.Background(), userFilter, pageSize, pageToken)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Empty(t, nextPageToken)
	assert.Equal(t, users[2], result[0])

	// Clean up test data
	_, err = collection.DeleteMany(context.Background(), bson.D{
		{"_id", bson.D{
			{"$in", []string{"user1", "user2", "user3"}},
		}},
	})
	require.NoError(t, err)
}

func TestMongoDB_ListUsers_SuccessWithFilter(t *testing.T) {
	collection := testMongoStorage.Database().Collection(UserCollection)
	users := []interface{}{
		storage.User{
			ID:        "user1",
			FirstName: "Alice",
			LastName:  "Smith",
			Email:     "alice@example.com",
			Nickname:  "alice",
			Country:   "USA",
			CreatedAt: testTimeForStorage(time.Now().Add(-4 * time.Hour)),
			UpdatedAt: testTimeForStorage(time.Now().Add(-3 * time.Hour)),
		},
		storage.User{
			ID:        "user2",
			FirstName: "Bob",
			LastName:  "Brown",
			Email:     "bob@example.com",
			Nickname:  "bobby",
			Country:   "USA",
			CreatedAt: testTimeForStorage(time.Now().Add(-3 * time.Hour)),
			UpdatedAt: testTimeForStorage(time.Now().Add(-2 * time.Hour)),
		},
		storage.User{
			ID:        "user3",
			FirstName: "Charlie",
			LastName:  "Davis",
			Email:     "charlie@example.com",
			Nickname:  "charlie",
			Country:   "UK",
			CreatedAt: testTimeForStorage(time.Now().Add(-2 * time.Hour)),
			UpdatedAt: testTimeForStorage(time.Now().Add(-1 * time.Hour)),
		},
	}
	_, err := collection.InsertMany(context.Background(), users)
	require.NoError(t, err)

	// Test ListUsers
	country := "USA"
	userFilter := storage.UserFilter{
		Country: &country,
	}
	pageSize := 1
	pageToken := ""

	result, nextPageToken, err := testMongoStorage.ListUsers(context.Background(), userFilter, pageSize, pageToken)
	require.NoError(t, err)
	assert.Len(t, result, pageSize)
	assert.NotEmpty(t, nextPageToken)
	assert.Equal(t, users[0], result[0])

	// Test ListUsers with next page token
	pageToken = nextPageToken
	result, nextPageToken, err = testMongoStorage.ListUsers(context.Background(), userFilter, pageSize, pageToken)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Empty(t, nextPageToken)
	assert.Equal(t, users[1], result[0])

	// Clean up test data
	_, err = collection.DeleteMany(context.Background(), bson.D{
		{"_id", bson.D{
			{"$in", []string{"user1", "user2", "user3"}},
		}},
	})
	require.NoError(t, err)
}

func TestMongoDB_ListUsers_InvalidPageTokenError(t *testing.T) {
	userFilter := storage.UserFilter{}
	pageSize := 1
	invalidPageToken := "invalid-token"

	_, _, err := testMongoStorage.ListUsers(context.Background(), userFilter, pageSize, invalidPageToken)
	assert.Error(t, err)
	var errCommon common.Error
	assert.ErrorAs(t, err, &errCommon)
	assert.Equal(t, common.ErrTypeInvalidArgument, errCommon.Type())
}

func TestMongoDB_ListUsers_SuccessInvalidPageSize(t *testing.T) {
	collection := testMongoStorage.Database().Collection(UserCollection)
	users := []interface{}{
		storage.User{
			ID:        "user1",
			FirstName: "Alice",
			LastName:  "Smith",
			Email:     "alice@example.com",
			Nickname:  "alice",
			Country:   "USA",
			CreatedAt: testTimeForStorage(time.Now().Add(-4 * time.Hour)),
			UpdatedAt: testTimeForStorage(time.Now().Add(-3 * time.Hour)),
		},
		storage.User{
			ID:        "user2",
			FirstName: "Bob",
			LastName:  "Brown",
			Email:     "bob@example.com",
			Nickname:  "bobby",
			Country:   "USA",
			CreatedAt: testTimeForStorage(time.Now().Add(-3 * time.Hour)),
			UpdatedAt: testTimeForStorage(time.Now().Add(-2 * time.Hour)),
		},
		storage.User{
			ID:        "user3",
			FirstName: "Charlie",
			LastName:  "Davis",
			Email:     "charlie@example.com",
			Nickname:  "charlie",
			Country:   "UK",
			CreatedAt: testTimeForStorage(time.Now().Add(-2 * time.Hour)),
			UpdatedAt: testTimeForStorage(time.Now().Add(-1 * time.Hour)),
		},
	}
	_, err := collection.InsertMany(context.Background(), users)
	require.NoError(t, err)

	// Test ListUsers
	userFilter := storage.UserFilter{}
	pageSize := 10000
	pageToken := ""

	result, nextPageToken, err := testMongoStorage.ListUsers(context.Background(), userFilter, pageSize, pageToken)
	require.NoError(t, err)
	assert.Len(t, result, len(users))
	assert.Empty(t, nextPageToken)
	assert.Equal(t, users[0], result[0])
	assert.Equal(t, users[1], result[1])
	assert.Equal(t, users[2], result[2])

	// Clean up test data
	_, err = collection.DeleteMany(context.Background(), bson.D{
		{"_id", bson.D{
			{"$in", []string{"user1", "user2", "user3"}},
		}},
	})
	require.NoError(t, err)
}
