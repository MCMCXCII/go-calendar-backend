package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

type Params struct {
	URL string
}

func New(ctx context.Context, p Params) (*DB, error) {
	if p.URL == "" {
		return nil, ErrDatabaseUrlEmpty
	}

	pool, err := pgxpool.New(ctx, p.URL)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	db := &DB{
		pool: pool,
	}

	if err := db.Ready(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return db, nil
}

func (db *DB) Ready(ctx context.Context) error {
	if db.pool == nil {
		return ErrDatabaseNotReady
	}

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := db.pool.Ping(pingCtx); err != nil {
		return fmt.Errorf("postgres ping: %w", err)
	}

	return nil
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}
