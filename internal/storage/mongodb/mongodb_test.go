package mongodb

import (
	"context"
	"fmt"
	"github.com/alenalato/users-service/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"log"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var testMongoStorage *MongoDB

const testDbName = "test"

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct test pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pull mongodb docker image
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "4.4",
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	var dbClient *mongo.Client

	// exponential backoff-retry
	err = pool.Retry(func() error {
		dbClient, err = newMongoDBClient(fmt.Sprintf(
			"mongodb://%s:%s",
			resource.Container.NetworkSettings.Gateway,
			resource.GetPort("27017/tcp"),
		))
		if err != nil {
			return err
		}
		return dbClient.Ping(context.TODO(), nil)
	})

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	testMongoStorage, err = NewMongoDB(dbClient, testDbName)

	defer func() {
		// kill and remove the container
		if err = pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
		// disconnect mongodb client
		if err = dbClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// run tests
	m.Run()
}

func TestMongoDB_UniqueIndexes(t *testing.T) {
	insertCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := testMongoStorage.GetDataBase().Collection(UserCollection)

	users := []interface{}{
		storage.UserDetails{
			ID:           "ff6e31dd-fa91-4c23-8cf8-c5c7f9b364c9",
			FirstName:    "John",
			LastName:     "Doe",
			Email:        "john@doe.com",
			Nickname:     "johndoe",
			PasswordHash: "password",
			Country:      "USA",
			CreatedAt:    time.Now().Add(-1 * time.Hour),
			UpdatedAt:    time.Now(),
		},
	}

	_, err := collection.InsertMany(
		insertCtx,
		users,
	)
	require.NoError(t, err)

	// Test unique index on ID
	_, err = testMongoStorage.database.Collection(UserCollection).
		InsertOne(insertCtx, storage.UserDetails{
			ID: "ff6e31dd-fa91-4c23-8cf8-c5c7f9b364c9",
		})

	assert.Error(t, err)
	assert.True(t, mongo.IsDuplicateKeyError(err))

	// Test unique index on Email
	_, err = testMongoStorage.database.Collection(UserCollection).
		InsertOne(insertCtx, storage.UserDetails{
			ID:    uuid.New().String(),
			Email: "john@doe.com",
		})

	assert.Error(t, err)
	assert.True(t, mongo.IsDuplicateKeyError(err))

	// Test unique index on Nickname
	_, err = testMongoStorage.database.Collection(UserCollection).
		InsertOne(insertCtx, storage.UserDetails{
			ID:       uuid.New().String(),
			Nickname: "johndoe",
		})

	assert.Error(t, err)
	assert.True(t, mongo.IsDuplicateKeyError(err))

	// Clean up test data
	_, err = collection.DeleteMany(
		context.Background(),
		bson.D{
			{Key: "_id", Value: "ff6e31dd-fa91-4c23-8cf8-c5c7f9b364c9"},
		},
	)
}
