package common

import "time"

type TimeProvider interface {
	Now() time.Time
}

type Time struct {
	time.Time
}

func (t Time) Now() time.Time {
	return time.Now()
}

func NewTime() Time {
	return Time{}
}
