package provider

import "time"

type TimeProvider interface {
	Now() time.Time
}

type timeProvider struct {
}

func (timeProvider) Now() time.Time {
	return time.Now()
}

func NewTimeProvider() TimeProvider {
	return timeProvider{}
}
