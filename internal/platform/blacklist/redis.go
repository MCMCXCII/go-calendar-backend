package blacklist

import (
	"context"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

type Params struct {
	Client *redis.Client
}

func New(p Params) *Redis {
	if p.Client == nil {
		panic("blacklist redis client is nill")
	}

	return &Redis{client: p.Client}
}

func (r *Redis) Revoke(ctx context.Context, tokenID string, ttl time.Duration) error {
	if tokenID == "" {
		return ErrTokenIDEmpty
	}

	if ttl <= 0 {
		return ErrInvalidTTL
	}

	key := revokedTokenKey(tokenID)

	if err := r.client.Set(ctx, key, "1", ttl).Err(); err != nil {
		return fmt.Errorf("save revoked token: %w", err)
	}

	return nil
}

func (r *Redis) IsRevoked(ctx context.Context, tokenID string) (bool, error) {
	if tokenID == "" {
		return false, ErrTokenIDEmpty
	}

	key := revokedTokenKey(tokenID)

	exist, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("check revoked token: %w", err)
	}

	return exist > 0, nil
}

func revokedTokenKey(tokenID string) string {
	return "revoked:access:" + tokenID
}
