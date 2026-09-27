package v1

import (
	"context"
	"errors"
	"net/http"
	"project/internal/auth/domain"
	"project/internal/auth/service"
	"project/pkg/render"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

func writeError(ctx context.Context, w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		render.Error(ctx, w, err, http.StatusUnauthorized, "invalid credentials")
	case errors.Is(err, domain.ErrEmailEmpty):
		render.Error(ctx, w, err, http.StatusBadRequest, "email is empty")
	case errors.Is(err, domain.ErrPasswordEmpty):
		render.Error(ctx, w, err, http.StatusBadRequest, "password is empty")
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		render.Error(ctx, w, err, http.StatusConflict, "email already exists")
	case errors.Is(err, service.ErrTokenExpired):
		render.Error(ctx, w, err, http.StatusUnauthorized, "token expired")
	case errors.Is(err, service.ErrTokenIDIsEmpty):
		render.Error(ctx, w, err, http.StatusUnauthorized, "token id is empty")
	default:
		render.Error(ctx, w, err, http.StatusInternalServerError, "internal error")
	}
}
