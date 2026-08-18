package v1

import (
	"encoding/json"
	"net/http"
	"project/internal/auth/service"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func (v *V1) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}
	if err := v.validate.Struct(req); err != nil {
		writeValidationError(w, err)
		return
	}

	result, err := v.uc.Login(r.Context(), service.LoginParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, LoginResponse{AccessToken: result.AccessToken})
}
