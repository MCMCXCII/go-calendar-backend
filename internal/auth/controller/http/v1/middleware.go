package v1

import (
	"context"
	"net/http"
	"strings"

	"project/pkg/render"
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
			ctx := r.Context()

			header := r.Header.Get("Authorization")
			tokenString, ok := strings.CutPrefix(header, "Bearer ")
			if !ok || tokenString == "" {
				render.Error(ctx, w, ErrUnauthorized, http.StatusUnauthorized, "missing authorization header")
				return
			}

			info, err := tokenParser.ParseAccessToken(tokenString)
			if err != nil {
				render.Error(ctx, w, ErrUnauthorized, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			revoked, err := blacklist.IsRevoked(ctx, info.TokenID)
			if err != nil {
				render.Error(ctx, w, err, http.StatusInternalServerError, "check blacklist failed")
				return
			}
			if revoked {
				render.Error(ctx, w, ErrUnauthorized, http.StatusUnauthorized, "token revoked")
				return
			}

			ctx = context.WithValue(ctx, tokenInfoContextKey{}, info)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
