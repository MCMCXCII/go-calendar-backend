package domain

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventTypeMeeting  EventType = "meeting"
	EventTypeTask     EventType = "task"
	EventTypeReminder EventType = "reminder"
	EventTypeOther    EventType = "other"
)

func (t EventType) IsValid() bool {
	switch t {
	case EventTypeMeeting, EventTypeTask, EventTypeReminder, EventTypeOther:
		return true
	default:
		return false
	}
}

type Event struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       string
	Type        EventType
	CustomType  string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
