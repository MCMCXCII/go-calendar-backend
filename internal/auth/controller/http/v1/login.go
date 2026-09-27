package v1

import (
	"encoding/json"
	"net/http"

	"project/internal/auth/service"
	"project/pkg/render"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func (v *V1) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(ctx, w, err, http.StatusBadRequest, "json decode error")
		return
	}

	if err := v.Validate(req); err != nil {
		writeError(ctx, w, err)
		return
	}

	token, err := v.uc.Login(ctx, service.LoginParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		writeError(ctx, w, err)
		return
	}

	render.JSON(w, LoginResponse{AccessToken: token.AccessToken}, http.StatusOK)
}
