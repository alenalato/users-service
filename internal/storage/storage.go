package storage

import "context"

type UserStorage interface {
	CreateUser(ctx context.Context, userDetails UserDetails) (*User, error)
	UpdateUser(ctx context.Context, userId string, userUpdate UserUpdate) (*User, error)
}
