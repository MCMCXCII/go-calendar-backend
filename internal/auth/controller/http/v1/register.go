package v1

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"project/internal/auth/service"
	"project/pkg/render"
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
	ctx := r.Context()

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(ctx, w, err, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := v.Validate(req); err != nil {
		writeError(ctx, w, err)
		return
	}

	registered, err := v.uc.Register(r.Context(), service.RegisterParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		writeError(ctx, w, err)
		return
	}

	render.JSON(w, RegisterResponse{
		UserID:  registered.UserID,
		Message: "registered",
	}, http.StatusCreated)
}
