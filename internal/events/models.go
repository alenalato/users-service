package events

import "time"

type UserEventType string

const EventTypeCreated UserEventType = "created"
const EventTypeUpdated UserEventType = "updated"
const EventTypeDeleted UserEventType = "deleted"

// UserEvent represents a user event
// It is used to emit events to the event bus
// EventType is one of the following: created, updated, deleted
// EventTime is the time when the event was emitted
// EventMask is a list of fields that were changed in the user for the updated event
type UserEvent struct {
	EventType UserEventType `json:"event_type"`
	EventTime time.Time     `json:"event_time"`

	UserId    string    `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	EventMask []string `json:"event_mask"`
}
