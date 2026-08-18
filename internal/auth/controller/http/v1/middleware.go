package v1

import (
	"context"
	"net/http"
	"strings"

	"project/pkg/token"
)

type tokenInfoContextKey struct{}

type TokenParser interface {
	ParseAccessToken(tokenString string) (token.Info, error)
}

type BlackList interface {
	IsRevoked(ctx context.Context, tokenID string) (bool, error)
}

func Auth(tokenParser TokenParser, blacklist BlackList) func(http.Handler) http.Handler {
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

			ctx := context.WithValue(r.Context(), tokenInfoContextKey{}, info)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
