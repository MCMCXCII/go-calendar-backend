package httpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"project/internal/events/domain"
	"project/internal/events/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type app interface {
	CreateEvent(ctx context.Context, p service.CreateEventParams) (service.CreateEventResult, error)
	GetEvent(ctx context.Context, userID, eventID uuid.UUID) (domain.Event, error)
	ListEvents(ctx context.Context, p service.ListEventsParams) ([]domain.Event, error)
	UpdateEvent(ctx context.Context, p service.UpdateEventParams) error
	DeleteEvent(ctx context.Context, userID, eventID uuid.UUID) error
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

func (s *Server) handleCreateEvent(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}
	if err := s.validate.Struct(req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	result, err := s.app.CreateEvent(r.Context(), service.CreateEventParams{
		UserID:      userID,
		Title:       req.Title,
		Type:        domain.EventType(req.Type),
		CustomType:  req.CustomType,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, CreateEventResponse{ID: result.EventID.String(), Message: "created"})
}

func (s *Server) handleGetEvent(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	eventID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid event id"})
		return
	}

	event, err := s.app.GetEvent(r.Context(), userID, eventID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toEventResponse(event))
}

func (s *Server) handleListEvents(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	q := r.URL.Query()

	events, err := s.app.ListEvents(r.Context(), service.ListEventsParams{
		UserID: userID,
		Day:    q.Get("day"),
		Week:   q.Get("week"),
		Month:  q.Get("month"),
		From:   q.Get("from"),
		To:     q.Get("to"),
	})
	if err != nil {
		writeError(w, err)
		return
	}

	resp := ListEventsResponse{Events: make([]EventResponse, 0, len(events))}
	for _, e := range events {
		resp.Events = append(resp.Events, toEventResponse(e))
	}

	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	eventID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid event id"})
		return
	}

	var req UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}
	if err := s.validate.Struct(req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err = s.app.UpdateEvent(r.Context(), service.UpdateEventParams{
		EventID:     eventID,
		UserID:      userID,
		Title:       req.Title,
		Type:        domain.EventType(req.Type),
		CustomType:  req.CustomType,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, MessageResponse{Message: "updated"})
}

func (s *Server) handleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	eventID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid event id"})
		return
	}

	if err := s.app.DeleteEvent(r.Context(), userID, eventID); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, MessageResponse{Message: "deleted"})
}
