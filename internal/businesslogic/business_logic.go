package businesslogic

import (
	"context"
)

//go:generate mockgen -destination=business_logic_mock.go -package=businesslogic github.com/alenalato/users-service/internal/businesslogic UserManager

type UserManager interface {
	CreateUser(ctx context.Context, userDetails UserDetails) (*User, error)
	UpdateUser(ctx context.Context, userId string, userUpdate UserUpdate) (*User, error)
	DeleteUser(ctx context.Context, userId string) error
	ListUsers(ctx context.Context, userFilter UserFilter, pageSize int, pageToken string) ([]User, string, error)
}

//go:generate mockgen -destination=password_manager_mock.go -package=businesslogic github.com/alenalato/users-service/internal/businesslogic PasswordManager

type PasswordManager interface {
	GeneratePasswordHash(ctx context.Context, passwordDetails *PasswordDetails) error
	VerifyPassword(ctx context.Context, password string, passwordDetails *PasswordDetails) error
}
