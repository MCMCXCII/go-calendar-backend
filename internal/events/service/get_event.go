package service

import (
	"context"
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
		return domain.Event{}, fmt.Errorf("store.GetEvent: %w", err)
	}
	return event, nil
}
