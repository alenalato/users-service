package businesslogic

import (
	"context"
)

//go:generate mockgen -destination=business_logic_mock.go -package=businesslogic github.com/alenalato/users-service/internal/businesslogic UserManager

// UserManager is an interface for business logic operations related to user management
type UserManager interface {
	CreateUser(ctx context.Context, userDetails UserDetails) (*User, error)
	UpdateUser(ctx context.Context, userId string, userUpdate UserUpdate) (*User, error)
	DeleteUser(ctx context.Context, userId string) error
	ListUsers(ctx context.Context, userFilter UserFilter, pageSize int, pageToken string) ([]User, string, error)
}

//go:generate mockgen -destination=password_manager_mock.go -package=businesslogic github.com/alenalato/users-service/internal/businesslogic PasswordManager

// PasswordManager is an interface for password management operations
type PasswordManager interface {
	GeneratePasswordHash(ctx context.Context, passwordDetails *PasswordDetails) error
	VerifyPassword(ctx context.Context, passwordDetails *PasswordDetails) error
}
