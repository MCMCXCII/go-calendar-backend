package render

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"project/pkg/logger"
)

type Err struct {
	Error string `json:"error"`
}

type Message struct {
	Message string `json:"message"`
}

func Error(ctx context.Context, w http.ResponseWriter, err error, status int, message string) {
	wrapped := fmt.Errorf("%s: %w", message, err)
	logger.SetError(ctx, wrapped)

	JSON(w, Err{Error: unpack(err).Error()}, status)
}

func unpack(err error) error {
	for {
		e := errors.Unwrap(err)
		if e == nil {
			break
		}

		err = e
	}

	return err
}
