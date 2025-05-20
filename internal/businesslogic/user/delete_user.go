package user

import (
	"context"
	"github.com/alenalato/users-service/internal/events"
	"github.com/alenalato/users-service/internal/logger"
)

func (l *Logic) DeleteUser(ctx context.Context, userId string) error {
	// delete user in storage
	deleteErr := l.userStorage.DeleteUser(ctx, userId)
	if deleteErr != nil {
		return deleteErr
	}

	// emit user event
	userEvent := events.UserEvent{
		UserId:    userId,
		EventType: events.EventTypeDeleted,
		EventTime: l.time.Now(),
	}
	emitErr := l.eventEmitter.EmitUserEvent(ctx, userEvent)
	if emitErr != nil {
		logger.Log.Warn("User deleted without event emission: %v", userEvent)
	}

	return nil
}
