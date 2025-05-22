package password

import (
	"context"
	"github.com/alenalato/users-service/internal/common"
	"testing"

	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/stretchr/testify/assert"
)

func TestGeneratePasswordHash(t *testing.T) {
	bcrypt := NewBcrypt()
	ctx := context.Background()

	passwordDetails := &businesslogic.PasswordDetails{
		Text: "securepassword",
	}

	err := bcrypt.GeneratePasswordHash(ctx, passwordDetails)
	assert.NoError(t, err)
	assert.NotEmpty(t, passwordDetails.Hash)
}

func TestVerifyPassword(t *testing.T) {
	bcrypt := NewBcrypt()
	ctx := context.Background()

	passwordDetails := &businesslogic.PasswordDetails{
		Text: "securepassword",
	}

	// Generate hash first
	err := bcrypt.GeneratePasswordHash(ctx, passwordDetails)
	assert.NoError(t, err)

	// Verify the password
	err = bcrypt.VerifyPassword(ctx, passwordDetails)
	assert.NoError(t, err)
}

func TestVerifyPassword_Invalid(t *testing.T) {
	bcrypt := NewBcrypt()
	ctx := context.Background()

	passwordDetails := &businesslogic.PasswordDetails{
		Text: "wrongpassword",
		Hash: "$2a$14$invalidhash", // Invalid hash to simulate an error
	}

	err := bcrypt.VerifyPassword(ctx, passwordDetails)
	assert.Error(t, err)
	var errCommon common.Error
	assert.ErrorAs(t, err, &errCommon)
	assert.Equal(t, common.ErrTypeInvalidArgument, errCommon.Type())
}
