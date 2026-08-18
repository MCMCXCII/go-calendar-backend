package v1

import (
	"net/http"

	"project/internal/auth/service"
	"project/pkg/token"
)

func (v *V1) Logout(w http.ResponseWriter, r *http.Request) {
	tokenInfo, ok := r.Context().Value(tokenInfoContextKey{}).(token.Info)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	if err := v.uc.Logout(r.Context(), service.LogoutParams{
		TokenID:   tokenInfo.TokenID,
		ExpiresAt: tokenInfo.ExpiresAt,
	}); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, MessageResponse{Message: "logged out"})
}
