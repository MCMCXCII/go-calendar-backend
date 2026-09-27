package service

import (
	"context"
	"fmt"
	"time"
)

type LogoutParams struct {
	TokenID   string
	ExpiresAt time.Time
}

func (s *Service) Logout(ctx context.Context, p LogoutParams) error {
	if p.TokenID == "" {
		return ErrTokenIDIsEmpty
	}

	ttl := time.Until(p.ExpiresAt)
	if ttl <= 0 {
		return ErrTokenExpired
	}

	if err := s.blacklist.Revoke(ctx, p.TokenID, ttl); err != nil {
		return fmt.Errorf("blacklist.Revoke: %w", err)
	}
	return nil
}
