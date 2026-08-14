package postgres

import "errors"

var (
	ErrDatabaseUrlEmpty = errors.New("database url is empty")
	ErrDatabaseNotReady = errors.New("database is not ready")
)
