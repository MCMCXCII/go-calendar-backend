package httpserver

import (
	"context"
	"net/http"
	"strings"
)

type tokenInfoContextKey struct{}

func Auth(tokens tokenParser, blacklist blackList) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "authorization header is required", http.StatusUnauthorized)
				return
			}

			const prefix = "Bearer "

			if !strings.HasPrefix(authHeader, prefix) {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, prefix)
			if tokenString == "" {
				http.Error(w, "token is empty", http.StatusUnauthorized)
				return
			}
			tokenInfo, err := tokens.ParseAccessToken(tokenString)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			revoked, err := blacklist.IsRevoked(r.Context(), tokenInfo.TokenID)
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			if revoked {
				http.Error(w, "token is revoked", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), tokenInfoContextKey{}, tokenInfo)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
