package v1

import (
	"net/http"

	"project/internal/events/service"
	"project/pkg/render"
)

type ListEventsResponse struct {
	Events []EventResponse `json:"events"`
}

func (v *V1) ListEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := userIDFromContext(ctx)
	if !ok {
		render.Error(ctx, w, ErrUnauthorized, http.StatusUnauthorized, "unauthorized")
		return
	}

	q := r.URL.Query()

	events, err := v.uc.ListEvents(ctx, service.ListEventsParams{
		UserID: userID,
		Day:    q.Get("day"),
		Week:   q.Get("week"),
		Month:  q.Get("month"),
		From:   q.Get("from"),
		To:     q.Get("to"),
	})
	if err != nil {
		writeError(ctx, w, err)
		return
	}

	resp := ListEventsResponse{Events: make([]EventResponse, 0, len(events))}
	for _, e := range events {
		resp.Events = append(resp.Events, toEventResponse(e))
	}

	render.JSON(w, resp, http.StatusOK)
}
