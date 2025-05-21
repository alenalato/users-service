package user

import (
	"context"
	"github.com/alenalato/users-service/internal/events"
	"github.com/alenalato/users-service/internal/logger"
)

func (l *Logic) DeleteUser(ctx context.Context, userId string) error {
	// delete user in storage
	errDelete := l.userStorage.DeleteUser(ctx, userId)
	if errDelete != nil {
		return errDelete
	}

	// emit user event
	userEvent := events.UserEvent{
		UserId:    userId,
		EventType: events.EventTypeDeleted,
		EventTime: l.time.Now().UTC(),
	}
	errEmit := l.eventEmitter.EmitUserEvent(ctx, userEvent)
	if errEmit != nil {
		logger.Log.Warn("User deleted without event emission: %v", userEvent)
	}

	return nil
}
