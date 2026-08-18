package v1

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

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

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("write json response", "error", err)
	}
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEventNotFound):
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "event not found"})

	case errors.Is(err, domain.ErrInvalidTimeRange),
		errors.Is(err, domain.ErrInvalidType),
		errors.Is(err, domain.ErrCustomTypeRequired),
		errors.Is(err, domain.ErrCustomTypeNotAllowed),
		errors.Is(err, domain.ErrTitleRequired),
		errors.Is(err, domain.ErrInvalidPeriod),
		errors.Is(err, domain.ErrPeriodRequired):
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})

	default:
		slog.Error("internal error", "error", err)
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
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
