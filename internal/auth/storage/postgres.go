package storage

import (
	"context"
	"errors"
	"fmt"
	"project/internal/auth/domain"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

type Params struct {
	URL string
}

func New(ctx context.Context, p Params) (*Store, error) {
	if p.URL == "" {
		return nil, ErrDatabaseUrlEmpty
	}

	pool, err := pgxpool.New(ctx, p.URL)
	if err != nil {
		return nil, fmt.Errorf("error create pool: %w", err)
	}

	store := &Store{
		pool: pool,
	}

	if err := store.Ready(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) Ready(ctx context.Context) error {
	if s.pool == nil {
		return ErrDatabaseNotReady
	}

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := s.pool.Ping(pingCtx); err != nil {
		return fmt.Errorf("database ping error: %w", err)
	}

	return nil
}

func (s *Store) Close() error {
	//реализовать
	return nil
}

func (s *Store) CreateUser(ctx context.Context, user domain.User) error {
	if _, err := s.pool.Exec(ctx,
		`INSERT INTO users(id, email, password_hash) VALUES($1, $2, $3)`,
		user.ID, user.Email, user.PasswordHash,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrEmailAlreadyExists
		}
		return fmt.Errorf("insert create: %w", err)
	}
	return nil
}

func (s *Store) GetUser(ctx context.Context, email string) (domain.User, error) {
	var user domain.User

	row := s.pool.QueryRow(
		ctx,
		`SELECT id, email, password_hash FROM users WHERE email = $1`,
		email,
	)

	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, fmt.Errorf("scan user: %w", err)
	}

	return user, nil
}
