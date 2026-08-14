package domain

import (
	"time"

	"github.com/google/uuid"
)

type ListEventsParams struct {
	UserID   uuid.UUID
	From     time.Time
	To       time.Time
	CacheKey string
}
