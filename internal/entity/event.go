package entity

import "time"

type Event struct {
	ID        *int
	Title     string
	Time      time.Time
	// Owner     string
	// IsDone    bool
}

func (e *Event) isActive() bool {
	return e.Time.Before(time.Now())
}