package events

import "time"

const (
	StreamEmailHigh = "stream:email:high"
	StreamEmailLow  = "stream:email:low"
)

const (
	PriorityHigh = "HIGH"
	PriorityLow  = "LOW"
)

type EmailEvent struct {
	EventID   string                 `json:"event_id"`
	Type      string                 `json:"type"`
	Priority  string                 `json:"priority"`
	Source    string                 `json:"source"`
	Recipient string                 `json:"recipient"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
}
