package db

import "time"

type Notification struct {
	ID              int
	ChatID          int64
	Message         string
	NextRun         time.Time
	IntervalSeconds int
}
