package logger

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ctxErrKey struct{}

func SetError(ctx context.Context, err error) {
	if errPtr, ok := ctx.Value(ctxErrKey{}).(*error); ok {
		*errPtr = err
	}
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var reqErr error
		ctx := context.WithValue(r.Context(), ctxErrKey{}, &reqErr)

		next.ServeHTTP(w, r.WithContext(ctx))

		fields := []zap.Field{
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Duration("duration", time.Since(start)),
			zap.String("remote_addr", r.RemoteAddr),
		}

		if reqErr != nil {
			Log.Error("http request", append(fields, zap.Error(reqErr))...)
			return
		}

		Log.Info("http request", fields...)
	})
}
