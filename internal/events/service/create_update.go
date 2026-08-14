package service

import (
	"context"
	"fmt"
	"project/internal/events/domain"
	"time"

	"github.com/google/uuid"
)

type CreateEventParams struct {
	UserID      uuid.UUID
	Title       string
	Type        domain.EventType
	CustomType  string
	Description string
	StartTime   time.Time
	EndTime     time.Time
}

type CreateEventResult struct {
	EventID uuid.UUID
}

func (s *Service) CreateEvent(ctx context.Context, p CreateEventParams) (CreateEventResult, error) {
	if err := validateEventInput(p.Title, p.Type, p.CustomType, p.StartTime, p.EndTime); err != nil {
		return CreateEventResult{}, err
	}

	event := domain.Event{
		ID:          uuid.New(),
		UserID:      p.UserID,
		Title:       p.Title,
		Type:        p.Type,
		CustomType:  p.CustomType,
		Description: p.Description,
		StartTime:   p.StartTime,
		EndTime:     p.EndTime,
	}

	if err := s.store.CreateEvent(ctx, event); err != nil {
		return CreateEventResult{}, fmt.Errorf("create event: %w", err)
	}

	return CreateEventResult{EventID: event.ID}, nil
}

type UpdateEventParams struct {
	EventID     uuid.UUID
	UserID      uuid.UUID
	Title       string
	Type        domain.EventType
	CustomType  string
	Description string
	StartTime   time.Time
	EndTime     time.Time
}

func (s *Service) UpdateEvent(ctx context.Context, p UpdateEventParams) error {
	if err := validateEventInput(p.Title, p.Type, p.CustomType, p.StartTime, p.EndTime); err != nil {
		return err
	}

	event := domain.Event{
		ID:          p.EventID,
		UserID:      p.UserID,
		Title:       p.Title,
		Type:        p.Type,
		CustomType:  p.CustomType,
		Description: p.Description,
		StartTime:   p.StartTime,
		EndTime:     p.EndTime,
	}

	if err := s.store.UpdateEvent(ctx, event); err != nil {
		return fmt.Errorf("update event: %w", err)
	}
	return nil
}

func validateEventInput(title string, t domain.EventType, customType string, start, end time.Time) error {
	if title == "" {
		return domain.ErrTitleRequired
	}
	if !t.IsValid() {
		return domain.ErrInvalidType
	}
	if t == domain.EventTypeOther && customType == "" {
		return domain.ErrCustomTypeRequired
	}
	if t != domain.EventTypeOther && customType != "" {
		return domain.ErrCustomTypeNotAllowed
	}
	if !start.Before(end) {
		return domain.ErrInvalidTimeRange
	}
	return nil
}
