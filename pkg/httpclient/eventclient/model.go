package eventsclient

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	CustomType  string    `json:"custom_type,omitempty"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EventInput struct {
	Title       string
	Type        string
	CustomType  string
	Description string
	StartTime   time.Time
	EndTime     time.Time
}

type ListFilter struct {
	Day   string
	Week  string
	Month string
	From  time.Time
	To    time.Time
}
