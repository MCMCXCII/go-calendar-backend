package storage

import (
	"context"
	"errors"
	"fmt"
	"project/internal/events/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

func New(p *pgxpool.Pool) *Store {
	return &Store{pool: p}
}

func (s *Store) CreateEvent(ctx context.Context, e domain.Event) error {
	if _, err := s.pool.Exec(ctx,
		`INSERT INTO events(id, user_id, title, type, custom_type, description, start_time, end_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		e.ID, e.UserID, e.Title, e.Type, nullableString(e.CustomType),
		nullableString(e.Description), e.StartTime, e.EndTime,
	); err != nil {
		return fmt.Errorf("insert event: %w", err)
	}
	return nil
}

func (s *Store) GetEvent(ctx context.Context, userID, eventID uuid.UUID) (domain.Event, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT id, user_id, title, type, custom_type, description, start_time, end_time, created_at, updated_at
		 FROM events WHERE id = $1 AND user_id = $2`,
		eventID, userID,
	)
	return scanEvent(row)
}

func (s *Store) UpdateEvent(ctx context.Context, e domain.Event) error {
	tag, err := s.pool.Exec(ctx,
		`UPDATE events
		 SET title = $1, type = $2, custom_type = $3, description = $4,
		     start_time = $5, end_time = $6, updated_at = now()
		 WHERE id = $7 AND user_id = $8`,
		e.Title, e.Type, nullableString(e.CustomType), nullableString(e.Description), e.StartTime, e.EndTime, e.ID, e.UserID,
	)
	if err != nil {
		return fmt.Errorf("update event: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrEventNotFound
	}
	return nil
}

func (s *Store) ListEvents(ctx context.Context, p domain.ListEventsParams) ([]domain.Event, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, user_id, title, type, custom_type, description, start_time, end_time, created_at, updated_at
		 FROM events
		 WHERE user_id = $1 AND start_time < $3 AND end_time > $2
		 ORDER BY start_time`,
		p.UserID, p.From, p.To,
	)
	if err != nil {
		return nil, fmt.Errorf("query events: %w", err)
	}
	defer rows.Close()

	events := make([]domain.Event, 0)
	for rows.Next() {
		e, err := scanEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}
	return events, nil
}

func (s *Store) DeleteEvent(ctx context.Context, userID, eventID uuid.UUID) error {
	tag, err := s.pool.Exec(ctx,
		`DELETE FROM events WHERE id = $1 AND user_id = $2`,
		eventID, userID,
	)
	if err != nil {
		return fmt.Errorf("delete event: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrEventNotFound
	}
	return nil
}

type row interface {
	Scan(dest ...any) error
}

func scanEvent(r row) (domain.Event, error) {
	var (
		e           domain.Event
		customType  *string
		description *string
	)

	if err := r.Scan(
		&e.ID, &e.UserID, &e.Title, &e.Type, &customType, &description,
		&e.StartTime, &e.EndTime, &e.CreatedAt, &e.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Event{}, domain.ErrEventNotFound
		}
		return domain.Event{}, fmt.Errorf("scan event: %w", err)
	}

	if customType != nil {
		e.CustomType = *customType
	}
	if description != nil {
		e.Description = *description
	}

	return e, nil
}

func nullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
