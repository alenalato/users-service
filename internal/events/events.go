package events

import "context"

//go:generate mockgen -destination=events_mock.go -package=events github.com/alenalato/users-service/internal/events EventEmitter

type EventEmitter interface {
	EmitUserEvent(ctx context.Context, userEvent UserEvent) error
}
