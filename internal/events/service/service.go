package service

import (
	"context"
	"project/internal/events/domain"

	"github.com/google/uuid"
)

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
