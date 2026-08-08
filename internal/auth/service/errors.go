package service

import "errors"

var (
	ErrEmailEmpty         = errors.New("email is empty")
	ErrPasswordEmpty      = errors.New("password is empty")
	ErrInvalidCredentials = errors.New("invalid email or password")

	ErrJWTSecretEmpty = errors.New("JWT secret is empty")
	ErrInvalidToken   = errors.New("invalid token")
)
