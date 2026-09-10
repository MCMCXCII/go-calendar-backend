package eventsclient

import "errors"

var (
	ErrValidation   = errors.New("validation error")
	ErrUnauthorized = errors.New("unauthorized")
	ErrNotFound     = errors.New("event not found")
	ErrInternal     = errors.New("internal server error")
)
