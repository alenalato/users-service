package businesslogic

import (
	"context"
)

type UserManager interface {
	CreateUser(ctx context.Context, userDetails UserDetails) (*User, error)
}

type PasswordManager interface {
	GeneratePasswordHash(ctx context.Context, passwordDetails *PasswordDetails) error
	VerifyPassword(ctx context.Context, password string, passwordDetails *PasswordDetails) error
}
