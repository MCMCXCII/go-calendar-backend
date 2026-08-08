package token

import "errors"

var (
	ErrJWTSecretEmpty = errors.New("JWT secret is empty")
	ErrInvalidToken   = errors.New("invalid token")
)
