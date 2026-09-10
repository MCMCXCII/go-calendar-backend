package authclient

import "errors"

var (
	ErrValidation         = errors.New("validation error")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInternal           = errors.New("internal server error")
)
