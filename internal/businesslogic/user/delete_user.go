package user

import (
	"context"
)

func (l *Logic) DeleteUser(ctx context.Context, userId string) error {
	// delete user in storage
	return l.userStorage.DeleteUser(ctx, userId)
}
