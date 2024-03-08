package msgbroker

import (
	"time"
)

type Log struct {
	Id          int
	ProjectId   int
	Name        string
	Description string
	Priority    int
	Removed     bool
	EventTime   time.Time
}
