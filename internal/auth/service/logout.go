package service

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type LogoutParams struct {
	TokenID   string
	ExpiresAt time.Time
}

func (s *Service) Logout(ctx context.Context, p LogoutParams) error {
	if p.TokenID == "" {
		return errors.New("token id is empty")
	}

	ttl := time.Until(p.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("invalid token")
	}

	if err := s.blacklist.Revoke(ctx, p.TokenID, ttl); err != nil {
		return fmt.Errorf("revoke access token: %w", err)
	}
	return nil
}
