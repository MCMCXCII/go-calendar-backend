package domain

import "errors"

var (
	ErrEventNotFound        = errors.New("event not found")
	ErrInvalidTimeRange     = errors.New("start_time must be before end_time")
	ErrInvalidType          = errors.New("unknown event type")
	ErrCustomTypeRequired   = errors.New("custom_type is required when type is other")
	ErrCustomTypeNotAllowed = errors.New("custom_type must be empty unless type is other")
	ErrTitleRequired        = errors.New("title is required")
	ErrInvalidPeriod        = errors.New("invalid day/week/month value")
	ErrPeriodRequired       = errors.New("one of day, week, month, or from/to is required")
)
