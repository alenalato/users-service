package common

import "time"

//go:generate mockgen -source=./time.go -destination=./time_mock.go -package=common TimeProvider

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
