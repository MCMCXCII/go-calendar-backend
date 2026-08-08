package blacklist

import (
	"errors"
)

var (
	ErrTokenIDEmpty = errors.New("token ID is empty")
	ErrInvalidTTL   = errors.New("token TTL must be positive")
)
