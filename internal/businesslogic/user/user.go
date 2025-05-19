package user

import (
	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/internal/storage"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type Logic struct {
	converter       modelConverter
	passwordManager businesslogic.PasswordManager
	userStorage     storage.UserStorage
}

var _ businesslogic.UserManager = new(Logic)

// NewLogic creates a new user Logic
func NewLogic(passwordManager businesslogic.PasswordManager, userStorage storage.UserStorage) *Logic {
	return &Logic{
		converter:       newStorageModelConverter(),
		passwordManager: passwordManager,
		userStorage:     userStorage,
	}
}
