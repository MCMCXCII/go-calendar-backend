package v1

import (
	"net/http"

	"project/internal/auth/service"
	"project/pkg/render"
	"project/pkg/token"
)

type MessageResponse struct {
	Message string `json:"message"`
}

func (v *V1) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tokenInfo, ok := r.Context().Value(tokenInfoContextKey{}).(token.Info)
	if !ok {
		render.Error(ctx, w, ErrUnauthorized, http.StatusUnauthorized, "r.Context().Value() error")
		return
	}

	if err := v.uc.Logout(r.Context(), service.LogoutParams{
		TokenID:   tokenInfo.TokenID,
		ExpiresAt: tokenInfo.ExpiresAt,
	}); err != nil {
		writeError(ctx, w, err)
		return
	}

	render.JSON(w, MessageResponse{Message: "logged out"}, http.StatusOK)
}
