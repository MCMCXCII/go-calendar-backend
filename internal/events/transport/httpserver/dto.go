package httpserver

import "time"

type CreateEventRequest struct {
	Title       string    `json:"title" validate:"required"`
	Type        string    `json:"type" validate:"required"`
	CustomType  string    `json:"custom_type"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required"`
}

type UpdateEventRequest struct {
	Title       string    `json:"title" validate:"required"`
	Type        string    `json:"type" validate:"required"`
	CustomType  string    `json:"custom_type"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required"`
}

type EventResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	CustomType  string    `json:"custom_type,omitempty"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateEventResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type ListEventsResponse struct {
	Events []EventResponse `json:"events"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
