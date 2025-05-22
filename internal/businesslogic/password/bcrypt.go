package password

import (
	"context"
	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

// Bcrypt is a struct that implements the PasswordManager interface
type Bcrypt struct {
}

var _ businesslogic.PasswordManager = new(Bcrypt)

// GeneratePasswordHash generates a password hash using bcrypt and stores it in the given PasswordDetails struct
func (s *Bcrypt) GeneratePasswordHash(_ context.Context, passwordDetails *businesslogic.PasswordDetails) error {
	bytes, bcryptErr := bcrypt.GenerateFromPassword([]byte(passwordDetails.Text), 14)
	if bcryptErr != nil {
		logger.Log.Errorf("bcrypt error: %v", bcryptErr)

		return common.NewError(bcryptErr, common.ErrTypeInternal)
	}

	passwordDetails.Hash = string(bytes)

	return nil
}

// VerifyPassword verifies a password against a hash using bcrypt
func (s *Bcrypt) VerifyPassword(
	_ context.Context,
	passwordDetails *businesslogic.PasswordDetails,
) error {
	bcryptErr := bcrypt.CompareHashAndPassword([]byte(passwordDetails.Hash), []byte(passwordDetails.Text))
	if bcryptErr != nil {
		logger.Log.Errorf("bcrypt error: %v", bcryptErr)

		return common.NewError(bcryptErr, common.ErrTypeInvalidArgument)
	}

	return nil
}

// NewBcrypt creates a new Bcrypt instance
func NewBcrypt() *Bcrypt {
	return &Bcrypt{}
}
