package storage

import "context"

const MaxPageSize = 100

//go:generate mockgen -destination=storage_mock.go -package=storage github.com/alenalato/users-service/internal/storage UserStorage

// UserStorage is the repository interface for user storage
type UserStorage interface {
	CreateUser(ctx context.Context, userDetails UserDetails) (*User, error)
	UpdateUser(ctx context.Context, userId string, userUpdate UserUpdate) (*User, error)
	DeleteUser(ctx context.Context, userId string) error
	ListUsers(ctx context.Context, userFilter UserFilter, pageSize int, pageToken string) ([]User, string, error)
}
