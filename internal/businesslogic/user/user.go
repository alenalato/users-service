package user

import (
	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/events"
	"github.com/alenalato/users-service/internal/storage"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type Logic struct {
	time            common.TimeProvider
	converter       modelConverter
	passwordManager businesslogic.PasswordManager
	userStorage     storage.UserStorage
	eventEmitter    events.EventEmitter
}

var _ businesslogic.UserManager = new(Logic)

// NewLogic creates a new user Logic
func NewLogic(
	passwordManager businesslogic.PasswordManager,
	userStorage storage.UserStorage,
	eventEmitter events.EventEmitter,
) *Logic {
	return &Logic{
		time:            common.NewTime(),
		converter:       newStorageModelConverter(),
		passwordManager: passwordManager,
		userStorage:     userStorage,
		eventEmitter:    eventEmitter,
	}
}
