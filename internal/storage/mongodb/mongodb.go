package mongodb

import (
	"context"
	"github.com/alenalato/users-service/internal/logger"
	"github.com/alenalato/users-service/internal/storage"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"os"
)

const UserCollection = "user"

type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
}

var _ storage.UserStorage = new(MongoDB)

// Close closes the MongoDB connection
func (m *MongoDB) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

// GetDataBase returns the inner MongoDB database
func (m *MongoDB) GetDataBase() *mongo.Database {
	return m.database
}

// NewMongoDB creates a new MongoDB storage. If client is nil, it creates a new client using the MONGODB_URI environment variable
func NewMongoDB(client *mongo.Client, databaseName string) (*MongoDB, error) {
	if client == nil {
		logger.Log.Debugf("Creating new MongoDB client with URI: %s", os.Getenv("MONGODB_URI"))
		var newClientErr error
		client, newClientErr = newMongoDBClient(os.Getenv("MONGODB_URI"))
		if newClientErr != nil {
			return nil, newClientErr
		}
	}
	database := client.Database(databaseName)

	// Create unique index for user email
	_, indexErr := database.Collection(UserCollection).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: map[string]interface{}{
				"email": 1,
			},
			Options: options.Index().SetName("email-unique").SetUnique(true),
		},
	)
	if indexErr != nil {
		return nil, indexErr
	}

	// Create unique index for user nickname
	_, indexErr = database.Collection(UserCollection).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: map[string]interface{}{
				"nickname": 1,
			},
			Options: options.Index().SetName("nickname-unique").SetUnique(true),
		},
	)
	if indexErr != nil {
		return nil, indexErr
	}

	return &MongoDB{
		client:   client,
		database: database,
	}, nil
}

func newMongoDBClient(uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)

	return mongo.Connect(clientOptions)
}
