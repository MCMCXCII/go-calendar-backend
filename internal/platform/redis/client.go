package redis

import (
	"context"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
}

type Params struct {
	URL string
}

func New(ctx context.Context, p Params) (*Client, error) {
	if p.URL == "" {
		return nil, ErrURLEmpty
	}

	options, err := redis.ParseURL(p.URL)
	if err != nil {
		return nil, fmt.Errorf("parse redis URL: %w", err)
	}

	client := redis.NewClient(options)

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) Client() *redis.Client {
	return c.client
}

func (c *Client) Ready(ctx context.Context) error {
	if c == nil || c.client == nil {
		return ErrRedisNotReady
	}

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := c.client.Ping(pingCtx).Err(); err != nil {
		return fmt.Errorf("ping redis: %w", err)
	}

	return nil
}

func (c *Client) Close() error {
	if c == nil || c.client == nil {
		return nil
	}

	return c.client.Close()
}
