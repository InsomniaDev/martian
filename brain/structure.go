package brain

import (
	"time"
)

// Brain main struct for brain package
type Brain struct {
	Placeholder     string
	AutomationEvent []Event
	CurrentEvent    Event
	LastEvent       Event
	automationTime  time.Time
	redisURL        string
	redisPort       string
}

// Event is the structure for all of the events
type Event struct {
	ID    int       `json:"id"`
	Type  string    `json:"type"`
	Value string    `json:"value"`
	Time  time.Time `json:"time"`
}
