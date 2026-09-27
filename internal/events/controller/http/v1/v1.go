package v1

import (
	"context"
	"errors"

	"project/internal/events/domain"
	"project/internal/events/service"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type usecase interface {
	CreateEvent(ctx context.Context, p service.CreateEventParams) (service.CreateEventResult, error)
	GetEvent(ctx context.Context, userID, eventID uuid.UUID) (domain.Event, error)
	ListEvents(ctx context.Context, p service.ListEventsParams) ([]domain.Event, error)
	UpdateEvent(ctx context.Context, p service.UpdateEventParams) error
	DeleteEvent(ctx context.Context, userID, eventID uuid.UUID) error
}

type V1 struct {
	uc       usecase
	validate *validator.Validate
}

func New(uc usecase) *V1 {
	return &V1{
		uc:       uc,
		validate: validator.New(),
	}
}

func (v *V1) Validate(req any) error {
	err := v.validate.Struct(req)
	if err == nil {
		return nil
	}

	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return err
	}

	return ValidationError(ve[0])
}

func ValidationError(fe validator.FieldError) error {
	switch fe.Field() {
	case "Title":
		return domain.ErrTitleRequired
	case "Type":
		return domain.ErrInvalidType
	case "StartTime", "EndTime":
		return domain.ErrInvalidTimeRange
	}
	return errors.New("validation failed: " + fe.Field())
}

func toEventResponse(e domain.Event) EventResponse {
	return EventResponse{
		ID:          e.ID.String(),
		Title:       e.Title,
		Type:        string(e.Type),
		CustomType:  e.CustomType,
		Description: e.Description,
		StartTime:   e.StartTime,
		EndTime:     e.EndTime,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
