package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"project/internal/events/domain"
	"project/internal/events/service"
	"project/pkg/render"
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
	ctx := r.Context()

	userID, ok := userIDFromContext(ctx)
	if !ok {
		render.Error(ctx, w, ErrUnauthorized, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(ctx, w, err, http.StatusBadRequest, "json decode error")
		return
	}

	if err := v.Validate(req); err != nil {
		writeError(ctx, w, err)
		return
	}

	result, err := v.uc.CreateEvent(ctx, service.CreateEventParams{
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

	render.JSON(w, CreateEventResponse{ID: result.EventID.String(), Message: "created"}, http.StatusCreated)
}
