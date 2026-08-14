package httpserver

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"project/internal/auth/domain"
	"project/internal/auth/service"

	"github.com/go-playground/validator/v10"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write json response: %v", err)
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
		log.Printf("internal error: %v", err)
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
}

func writeValidationError(w http.ResponseWriter, err error) {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request"})
		return
	}
	if len(validationErrors) == 0 {
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
