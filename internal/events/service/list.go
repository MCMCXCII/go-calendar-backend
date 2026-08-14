package service

import (
	"context"
	"fmt"
	"project/internal/events/domain"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ListEventsParams struct {
	UserID uuid.UUID
	Day    string // "2025-04-15"
	Week   string // "2025-W16"
	Month  string // "2025-04"
	From   string // RFC3339, напр. "2025-04-15T10:00:00Z"
	To     string // RFC3339
}

func (s *Service) ListEvents(ctx context.Context, p ListEventsParams) ([]domain.Event, error) {
	q, err := buildListQuery(p)
	if err != nil {
		return nil, err
	}

	events, err := s.store.ListEvents(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}
	return events, nil
}

func buildListQuery(p ListEventsParams) (domain.ListEventsParams, error) {
	q := domain.ListEventsParams{UserID: p.UserID}

	switch {
	case p.From != "" && p.To != "":
		from, err := time.Parse(time.RFC3339, p.From)
		if err != nil {
			return domain.ListEventsParams{}, fmt.Errorf("%w: invalid from", domain.ErrInvalidPeriod)
		}
		to, err := time.Parse(time.RFC3339, p.To)
		if err != nil {
			return domain.ListEventsParams{}, fmt.Errorf("%w: invalid to", domain.ErrInvalidPeriod)
		}
		if !from.Before(to) {
			return domain.ListEventsParams{}, domain.ErrInvalidTimeRange
		}
		q.From, q.To = from, to
		q.CacheKey = ""

	case p.Day != "":
		from, err := time.Parse("2006-01-02", p.Day)
		if err != nil {
			return domain.ListEventsParams{}, fmt.Errorf("%w: invalid day", domain.ErrInvalidPeriod)
		}
		q.From = from
		q.To = from.AddDate(0, 0, 1)
		q.CacheKey = "day:" + p.Day

	case p.Week != "":
		from, to, err := parseISOWeek(p.Week)
		if err != nil {
			return domain.ListEventsParams{}, err
		}
		q.From, q.To = from, to
		q.CacheKey = "week:" + p.Week

	case p.Month != "":
		from, err := time.Parse("2006-01", p.Month)
		if err != nil {
			return domain.ListEventsParams{}, fmt.Errorf("%w: invalid month", domain.ErrInvalidPeriod)
		}
		q.From = from
		q.To = from.AddDate(0, 1, 0)
		q.CacheKey = "month:" + p.Month

	default:
		return domain.ListEventsParams{}, domain.ErrPeriodRequired
	}

	return q, nil
}

func parseISOWeek(s string) (time.Time, time.Time, error) {
	parts := strings.SplitN(s, "-W", 2)
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("%w: expected format YYYY-Www", domain.ErrInvalidPeriod)
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("%w: invalid year", domain.ErrInvalidPeriod)
	}

	week, err := strconv.Atoi(parts[1])
	if err != nil || week < 1 || week > 53 {
		return time.Time{}, time.Time{}, fmt.Errorf("%w: invalid week number", domain.ErrInvalidPeriod)
	}

	jan4 := time.Date(year, time.January, 4, 0, 0, 0, 0, time.UTC)

	weekday := int(jan4.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	week1Monday := jan4.AddDate(0, 0, -(weekday - 1))

	from := week1Monday.AddDate(0, 0, (week-1)*7)
	to := from.AddDate(0, 0, 7)

	return from, to, nil
}
