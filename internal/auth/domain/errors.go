package domain

import "errors"

var (
	ErrEmailEmpty         = errors.New("email is empty")
	ErrPasswordEmpty      = errors.New("password is empty")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrPasswordShort      = errors.New("password too short")
	ErrPasswordLong       = errors.New("password too long")
)
