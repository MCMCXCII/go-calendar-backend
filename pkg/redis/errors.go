package redis

import "errors"

var (
	ErrURLEmpty      = errors.New("redis URL is empty")
	ErrRedisNotReady = errors.New("redis client is not ready")
)
