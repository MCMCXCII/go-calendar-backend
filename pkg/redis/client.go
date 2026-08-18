package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Config struct {
	URL string `envconfig:"REDIS_URL" required:"true"`
}

type Client struct {
	Client *goredis.Client
}

func New(ctx context.Context, c Config) (*Client, error) {
	if c.URL == "" {
		return nil, ErrURLEmpty
	}

	options, err := goredis.ParseURL(c.URL)
	if err != nil {
		return nil, fmt.Errorf("parse redis URL: %w", err)
	}

	client := goredis.NewClient(options)

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &Client{Client: client}, nil
}

func (c *Client) Close() error {
	if c == nil || c.Client == nil {
		return nil
	}
	return c.Client.Close()
}
