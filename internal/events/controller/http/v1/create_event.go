package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"project/internal/events/domain"
	"project/internal/events/service"
)

type CreateEventRequest struct {
	Title       string    `json:"title" validate:"required"`
	Type        string    `json:"type" validate:"required"`
	CustomType  string    `json:"custom_type"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required"`
}

type CreateEventResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func (v *V1) CreateEvent(w http.ResponseWriter, r *http.Request) {
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
	if err := v.validate.Struct(req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	result, err := v.uc.CreateEvent(r.Context(), service.CreateEventParams{
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
