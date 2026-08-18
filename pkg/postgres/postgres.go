package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	URL string `envconfig:"DB_URL" required:"true"`
}

type Pool struct {
	*pgxpool.Pool
}

func New(ctx context.Context, c Config) (*Pool, error) {
	pool, err := pgxpool.New(ctx, c.URL)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pgxpool.Ping: %w", err)
	}

	return &Pool{Pool: pool}, nil
}

func (p *Pool) Close() {
	p.Pool.Close()
}
