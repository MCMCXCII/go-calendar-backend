package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"project/internal/events/domain"
	"project/internal/events/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UpdateEventRequest struct {
	Title       string    `json:"title" validate:"required"`
	Type        string    `json:"type" validate:"required"`
	CustomType  string    `json:"custom_type"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required"`
}

func (v *V1) UpdateEvent(w http.ResponseWriter, r *http.Request) {
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
	if err := v.validate.Struct(req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err = v.uc.UpdateEvent(r.Context(), service.UpdateEventParams{
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
