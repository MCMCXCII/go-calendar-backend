package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"project/internal/auth/domain"
	"project/internal/auth/service"
	"project/internal/platform/token"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type RegisterResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Message string    `json:"message"`
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := s.validate.Struct(req); err != nil {
		writeValidationError(w, err)
		return
	}

	result, err := s.app.Register(r.Context(), service.RegisterParams{
		Email: req.Email, Password: req.Password},
	)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailEmpty):
			http.Error(w, "email is required", http.StatusBadRequest)

		case errors.Is(err, service.ErrPasswordEmpty):
			http.Error(w, "password is required", http.StatusBadRequest)

		case errors.Is(err, domain.ErrEmailAlreadyExists):
			http.Error(w, "email already exists", http.StatusConflict)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(RegisterResponse{
		UserID:  result.UserID,
		Message: "registered",
	}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func writeValidationError(w http.ResponseWriter, err error) {
	var validationErrors validator.ValidationErrors

	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			switch e.Field() {
			case "password":
				switch e.Tag() {
				case "required":
					http.Error(w, "password is required", http.StatusBadRequest)
				case "min":
					http.Error(w, "password must be at least 8 characters", http.StatusBadRequest)
				case "max":
					http.Error(w, "password is too long", http.StatusBadRequest)
				}
			case "email":
				switch e.Tag() {
				case "required":
					http.Error(w, "email is required", http.StatusBadRequest)
				}
			}
		}
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := s.validate.Struct(req); err != nil {
		writeValidationError(w, err)
		return
	}

	result, err := s.app.Login(r.Context(), service.LoginParams{
		Email: req.Email, Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			http.Error(w, "invalid email or password", http.StatusUnauthorized)

		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(LoginResponse{
		AccessToken: result.AccessToken,
	}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	tokenInfo, ok := r.Context().Value(tokenInfoContextKey{}).(token.Info)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err := s.app.Logout(r.Context(), service.LogoutParams{TokenID: tokenInfo.TokenID,
		ExpiresAt: tokenInfo.ExpiresAt})
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
