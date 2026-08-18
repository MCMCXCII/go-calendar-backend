package v1

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"project/internal/auth/domain"
	"project/internal/auth/service"

	"github.com/go-playground/validator/v10"
)

type app interface {
	Register(ctx context.Context, req service.RegisterParams) (service.RegisterResult, error)
	Login(ctx context.Context, req service.LoginParams) (service.LoginResult, error)
	Logout(ctx context.Context, req service.LogoutParams) error
}

type V1 struct {
	uc       app
	validate *validator.Validate
}

func New(uc app) *V1 {
	return &V1{
		uc:       uc,
		validate: validator.New(),
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("write json response", "error", err)
	}
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrEmailEmpty):
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "email is required"})

	case errors.Is(err, service.ErrPasswordEmpty):
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "password is required"})

	case errors.Is(err, domain.ErrEmailAlreadyExists):
		writeJSON(w, http.StatusConflict, ErrorResponse{Error: "email already exists"})

	case errors.Is(err, service.ErrInvalidCredentials):
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "invalid email or password"})

	case errors.Is(err, service.ErrTokenExpired):
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "token expired"})

	default:
		slog.Error("internal error", "error", err)
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
}

func writeValidationError(w http.ResponseWriter, err error) {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) || len(validationErrors) == 0 {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request"})
		return
	}
	writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: validationMessage(validationErrors[0])})
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Field() {
	case "Email":
		switch fe.Tag() {
		case "required":
			return "email is required"
		case "email":
			return "email is invalid"
		}
	case "Password":
		switch fe.Tag() {
		case "required":
			return "password is required"
		case "min":
			return "password must be at least 8 characters"
		case "max":
			return "password is too long"
		}
	}
	return "invalid request"
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
