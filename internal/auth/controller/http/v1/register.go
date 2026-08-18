package v1

import (
	"encoding/json"
	"net/http"
	"project/internal/auth/service"

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

func (v *V1) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}
	if err := v.validate.Struct(req); err != nil {
		writeValidationError(w, err)
		return
	}

	result, err := v.uc.Register(r.Context(), service.RegisterParams{
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
