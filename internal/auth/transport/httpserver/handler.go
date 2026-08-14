package httpserver

import (
	"encoding/json"
	"net/http"

	"project/internal/auth/service"
	"project/internal/platform/token"
)

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}
	if err := s.validate.Struct(req); err != nil {
		writeValidationError(w, err)
		return
	}

	result, err := s.app.Register(r.Context(), service.RegisterParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, RegisterResponse{
		UserID:  result.UserID,
		Message: "registered",
	})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}
	if err := s.validate.Struct(req); err != nil {
		writeValidationError(w, err)
		return
	}

	result, err := s.app.Login(r.Context(), service.LoginParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, LoginResponse{AccessToken: result.AccessToken})
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	tokenInfo, ok := r.Context().Value(tokenInfoContextKey{}).(token.Info)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	if err := s.app.Logout(r.Context(), service.LogoutParams{
		TokenID:   tokenInfo.TokenID,
		ExpiresAt: tokenInfo.ExpiresAt,
	}); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, MessageResponse{Message: "logged out"})
}
