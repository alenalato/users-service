package user

import (
	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/events"
	"github.com/alenalato/users-service/internal/storage"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

// Logic is a struct that implements the UserManager interface
type Logic struct {
	// time is a time provider used for generating timestamps
	time common.TimeProvider
	// converter is a model converter used for converting between different models
	converter modelConverter
	// passwordManager is a password manager used for hashing and verifying passwords
	passwordManager businesslogic.PasswordManager
	// userStorage is a user storage used for storing and retrieving user data
	userStorage storage.UserStorage
	// eventEmitter is an event emitter used for emitting user events
	eventEmitter events.EventEmitter
}

var _ businesslogic.UserManager = new(Logic)

// NewLogic creates a new Logic instance
func NewLogic(
	passwordManager businesslogic.PasswordManager,
	userStorage storage.UserStorage,
	eventEmitter events.EventEmitter,
) *Logic {
	return &Logic{
		time:            common.NewTime(),
		converter:       newBusinessLogicModelConverter(),
		passwordManager: passwordManager,
		userStorage:     userStorage,
		eventEmitter:    eventEmitter,
	}
}
