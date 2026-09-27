package service

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrTokenExpired       = errors.New("token is already expired")
	ErrTokenIDIsEmpty     = errors.New("token id is empty")
)
