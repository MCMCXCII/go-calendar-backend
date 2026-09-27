package v1

import (
	"context"
	"net/http"
	"strings"

	"project/pkg/logger"
	"project/pkg/render"
	"project/pkg/token"

	"github.com/google/uuid"
	"go.uber.org/zap"
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

// Auth валидирует JWT самостоятельно (подпись + чёрный список в Redis),
// без похода в auth-service.
func Auth(tokenParser tokenParser, blacklist blackList) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			header := r.Header.Get("Authorization")
			logger.Log.Info("debug auth header", zap.String("raw_header", header), zap.Int("length", len(header)))

			tokenString, ok := strings.CutPrefix(header, "Bearer ")
			logger.Log.Info("debug after cut", zap.String("token", tokenString), zap.Bool("ok", ok))

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

			ctx = context.WithValue(ctx, userIDContextKey, info.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
