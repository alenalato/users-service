package storage

import "context"

const MaxPageSize = 10

type UserStorage interface {
	CreateUser(ctx context.Context, userDetails UserDetails) (*User, error)
	UpdateUser(ctx context.Context, userId string, userUpdate UserUpdate) (*User, error)
	DeleteUser(ctx context.Context, userId string) error
	ListUsers(ctx context.Context, userFilter UserFilter, pageSize int, pageToken string) ([]User, string, error)
}
