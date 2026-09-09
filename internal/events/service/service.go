package service

import (
	"context"
	"project/internal/events/domain"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=service.go -destination=./mocks/mocks.go -package=mocks

type store interface {
	CreateEvent(ctx context.Context, e domain.Event) error
	GetEvent(ctx context.Context, userID, eventID uuid.UUID) (domain.Event, error)
	ListEvents(ctx context.Context, p domain.ListEventsParams) ([]domain.Event, error)
	UpdateEvent(ctx context.Context, e domain.Event) error
	DeleteEvent(ctx context.Context, userID, eventID uuid.UUID) error
}

type Service struct {
	store store
}

type Params struct {
	Store store
}

func New(p Params) *Service {
	return &Service{store: p.Store}
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
