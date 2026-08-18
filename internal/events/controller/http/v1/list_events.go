package v1

import (
	"net/http"

	"project/internal/events/service"
)

type ListEventsResponse struct {
	Events []EventResponse `json:"events"`
}

func (v *V1) ListEvents(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	q := r.URL.Query()

	events, err := v.uc.ListEvents(r.Context(), service.ListEventsParams{
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
