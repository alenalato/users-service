package common

import "time"

//go:generate mockgen -source=./time.go -destination=./time_mock.go -package=common TimeProvider

// TimeProvider is an interface that wraps the time package
type TimeProvider interface {
	Now() time.Time
}

type Time struct {
	time.Time
}

// Now returns the current time
func (t Time) Now() time.Time {
	return time.Now()
}

// NewTime returns a new Time instance
func NewTime() Time {
	return Time{}
}
