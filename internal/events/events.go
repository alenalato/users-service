package events

import "context"

type EventEmitter interface {
	EmitUserEvent(ctx context.Context, userEvent UserEvent) error
}
