package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"project/internal/events/domain"
	"project/internal/events/service"
	"project/pkg/render"

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
	ctx := r.Context()

	userID, ok := userIDFromContext(ctx)
	if !ok {
		render.Error(ctx, w, ErrUnauthorized, http.StatusUnauthorized, "unauthorized")
		return
	}

	eventID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		render.Error(ctx, w, err, http.StatusBadRequest, "invalid event id")
		return
	}

	var req UpdateEventRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(ctx, w, err, http.StatusBadRequest, "json decode error")
		return
	}

	if err = v.Validate(req); err != nil {
		writeError(ctx, w, err)
		return
	}

	err = v.uc.UpdateEvent(ctx, service.UpdateEventParams{
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
		writeError(ctx, w, err)
		return
	}

	render.JSON(w, render.Message{Message: "updated"}, http.StatusOK)
}
