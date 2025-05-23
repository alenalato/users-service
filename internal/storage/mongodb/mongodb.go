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

// MongoDB is the MongoDB storage implementation of the UserStorage interface
type MongoDB struct {
	// client is the internal MongoDB client
	client *mongo.Client
	// database is the internal MongoDB database
	database *mongo.Database
}

var _ storage.UserStorage = new(MongoDB)

// Close closes the MongoDB connection for proper resource cleanup
func (m *MongoDB) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

// Database returns the inner MongoDB database
func (m *MongoDB) Database() *mongo.Database {
	return m.database
}

// NewMongoDB creates a new MongoDB storage.
// If client is nil, it creates a new client using the MONGODB_URI environment variable and connects to the database with the given name.
// It also creates unique indexes for user email and nickname.
func NewMongoDB(client *mongo.Client, databaseName string) (*MongoDB, error) {
	if client == nil {
		logger.Log.Debugf("Creating new MongoDB client with URI: %s", os.Getenv("MONGODB_URI"))
		var newClientErr error
		client, newClientErr = NewMongoDBClient(os.Getenv("MONGODB_URI"))
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

	// Create index for user created_at
	_, indexErr = database.Collection(UserCollection).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: map[string]interface{}{
				"created_at": 1,
			},
			Options: options.Index().SetName("created-at").SetUnique(false),
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

func NewMongoDBClient(uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)

	return mongo.Connect(clientOptions)
}
