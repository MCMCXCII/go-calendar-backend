package v1

import (
	"context"
	"errors"
	"net/http"
	"project/internal/events/domain"
	"project/pkg/render"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

func writeError(ctx context.Context, w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEventNotFound):
		render.Error(ctx, w, err, http.StatusNotFound, "event not found")

	case errors.Is(err, domain.ErrInvalidTimeRange):
		render.Error(ctx, w, err, http.StatusBadRequest, "invalid time range")
	case errors.Is(err, domain.ErrInvalidType):
		render.Error(ctx, w, err, http.StatusBadRequest, "invalid event type")
	case errors.Is(err, domain.ErrCustomTypeRequired):
		render.Error(ctx, w, err, http.StatusBadRequest, "custom type required")
	case errors.Is(err, domain.ErrCustomTypeNotAllowed):
		render.Error(ctx, w, err, http.StatusBadRequest, "custom type not allowed")
	case errors.Is(err, domain.ErrTitleRequired):
		render.Error(ctx, w, err, http.StatusBadRequest, "title required")
	case errors.Is(err, domain.ErrInvalidPeriod):
		render.Error(ctx, w, err, http.StatusBadRequest, "invalid period")
	case errors.Is(err, domain.ErrPeriodRequired):
		render.Error(ctx, w, err, http.StatusBadRequest, "period required")

	default:
		render.Error(ctx, w, err, http.StatusInternalServerError, "internal error")
	}
}
