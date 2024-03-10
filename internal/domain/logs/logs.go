package logs

import (
	"time"
)

type Log struct {
	Id          int       `json:"id"`
	ProjectId   int       `json:"project_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	EventTime   time.Time `json:"event_time"`
}
