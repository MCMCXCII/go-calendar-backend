package v1

import (
	"context"
	"net/http"
	"strings"

	"project/pkg/token"

	"github.com/google/uuid"
)

type contextKey string

const userIDContextKey contextKey = "userID"

func userIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDContextKey).(uuid.UUID)
	return id, ok
}

type tokenParser interface {
	ParseAccessToken(tokenString string) (token.Info, error)
}

type blackList interface {
	IsRevoked(ctx context.Context, tokenID string) (bool, error)
}

func Auth(tokenParser tokenParser, blacklist blackList) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			tokenString, ok := strings.CutPrefix(header, "Bearer ")
			if !ok || tokenString == "" {
				writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "missing or malformed authorization header"})
				return
			}

			info, err := tokenParser.ParseAccessToken(tokenString)
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "invalid or expired token"})
				return
			}

			revoked, err := blacklist.IsRevoked(r.Context(), info.TokenID)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
				return
			}
			if revoked {
				writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "token revoked"})
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, info.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
