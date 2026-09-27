package service

import (
	"context"
	"errors"
	"fmt"
	"project/internal/events/domain"

	"github.com/google/uuid"
)

func (s *Service) DeleteEvent(ctx context.Context, userID, eventID uuid.UUID) error {
	if err := s.store.DeleteEvent(ctx, userID, eventID); err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			return err
		}
		return fmt.Errorf("store.DeleteEvent: %w", err)
	}
	return nil
}
