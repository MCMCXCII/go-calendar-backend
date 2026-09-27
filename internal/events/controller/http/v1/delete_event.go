package v1

import (
	"net/http"

	"project/pkg/render"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (v *V1) DeleteEvent(w http.ResponseWriter, r *http.Request) {
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

	if err := v.uc.DeleteEvent(ctx, userID, eventID); err != nil {
		writeError(ctx, w, err)
		return
	}

	render.JSON(w, render.Message{Message: "deleted"}, http.StatusOK)
}
