package entity

import "time"

type Event struct {
	ID         *int
	Title      string
	Time       time.Time
	Status     string     // "pending" or "sent"
	ChatId     int
}
