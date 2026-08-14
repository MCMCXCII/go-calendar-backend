package service

import (
	"context"
	"errors"
	"fmt"
	"project/internal/events/domain"

	"github.com/google/uuid"
)

func (s *Service) GetEvent(ctx context.Context, userID, eventID uuid.UUID) (domain.Event, error) {
	event, err := s.store.GetEvent(ctx, userID, eventID)
	if err != nil {
		if err == domain.ErrEventNotFound {
			return domain.Event{}, err
		}
		return domain.Event{}, fmt.Errorf("get event: %w", err)
	}
	return event, nil
}

func (s *Service) DeleteEvent(ctx context.Context, userID, eventID uuid.UUID) error {
	if err := s.store.DeleteEvent(ctx, userID, eventID); err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			return err
		}
		return fmt.Errorf("delete event: %w", err)
	}
	return nil
}
