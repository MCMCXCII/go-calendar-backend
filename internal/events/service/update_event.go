package service

import (
	"context"
	"fmt"
	"project/internal/events/domain"
	"time"

	"github.com/google/uuid"
)

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
		return fmt.Errorf("store.UpdateEvent: %w", err)
	}
	return nil
}
