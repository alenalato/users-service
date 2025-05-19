package password

import (
	"context"
	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/internal/common"
	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct {
}

func (s *Bcrypt) GeneratePasswordHash(_ context.Context, passwordDetails *businesslogic.PasswordDetails) error {
	bytes, bcryptErr := bcrypt.GenerateFromPassword([]byte(passwordDetails.Text), 14)
	if bcryptErr != nil {
		return common.NewError(bcryptErr, common.ErrTypeInternal)
	}

	passwordDetails.Hash = string(bytes)

	return nil
}

func (s *Bcrypt) VerifyPassword(
	_ context.Context,
	password string,
	passwordDetails *businesslogic.PasswordDetails,
) error {
	bcryptErr := bcrypt.CompareHashAndPassword([]byte(passwordDetails.Hash), []byte(password))
	if bcryptErr != nil {
		return common.NewError(bcryptErr, common.ErrTypeInvalidArgument)
	}

	return nil
}

// NewBcrypt creates a new Bcrypt instance
func NewBcrypt() *Bcrypt {
	return &Bcrypt{}
}
