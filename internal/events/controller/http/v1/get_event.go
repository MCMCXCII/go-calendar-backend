package v1

import (
	"net/http"
	"time"

	"project/pkg/render"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type EventResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	CustomType  string    `json:"custom_type,omitempty"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (v *V1) GetEvent(w http.ResponseWriter, r *http.Request) {
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

	event, err := v.uc.GetEvent(ctx, userID, eventID)
	if err != nil {
		writeError(ctx, w, err)
		return
	}

	render.JSON(w, toEventResponse(event), http.StatusOK)
}
