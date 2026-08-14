package service

import "errors"

var (
	ErrEmailEmpty         = errors.New("email is empty")
	ErrPasswordEmpty      = errors.New("password is empty")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrTokenExpired       = errors.New("token is already expired")
)
