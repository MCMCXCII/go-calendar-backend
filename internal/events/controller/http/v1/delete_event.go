package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (v *V1) DeleteEvent(w http.ResponseWriter, r *http.Request) {
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

	if err := v.uc.DeleteEvent(r.Context(), userID, eventID); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, MessageResponse{Message: "deleted"})
}
