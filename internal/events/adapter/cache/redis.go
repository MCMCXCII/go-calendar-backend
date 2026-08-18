package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"project/internal/events/domain"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

type Config struct {
	TTL time.Duration `envconfig:"CACHE_TTL" default:"300s"`
}

type EventStore interface {
	CreateEvent(ctx context.Context, e domain.Event) error
	GetEvent(ctx context.Context, userID, eventID uuid.UUID) (domain.Event, error)
	ListEvents(ctx context.Context, p domain.ListEventsParams) ([]domain.Event, error)
	UpdateEvent(ctx context.Context, e domain.Event) error
	DeleteEvent(ctx context.Context, userID, eventID uuid.UUID) error
}

type CachingStore struct {
	next   EventStore
	client *goredis.Client
	ttl    time.Duration
}

type Params struct {
	Next   EventStore
	Client *goredis.Client
	Config Config
}

func New(p Params) *CachingStore {
	return &CachingStore{
		next:   p.Next,
		client: p.Client,
		ttl:    p.Config.TTL,
	}
}

func (c *CachingStore) CreateEvent(ctx context.Context, e domain.Event) error {
	if err := c.next.CreateEvent(ctx, e); err != nil {
		return err
	}
	c.invalidateUser(ctx, e.UserID)
	return nil
}

func (c *CachingStore) GetEvent(ctx context.Context, userID, eventID uuid.UUID) (domain.Event, error) {
	return c.next.GetEvent(ctx, userID, eventID)
}

func (c *CachingStore) ListEvents(ctx context.Context, p domain.ListEventsParams) ([]domain.Event, error) {
	if p.CacheKey == "" {
		return c.next.ListEvents(ctx, p)
	}

	key := cacheKey(p.UserID, p.CacheKey)

	if events, ok := c.getCached(ctx, key); ok {
		return events, nil
	}

	events, err := c.next.ListEvents(ctx, p)
	if err != nil {
		return nil, err
	}

	c.setCached(ctx, key, events)
	return events, nil
}

func (c *CachingStore) UpdateEvent(ctx context.Context, e domain.Event) error {
	if err := c.next.UpdateEvent(ctx, e); err != nil {
		return err
	}
	c.invalidateUser(ctx, e.UserID)
	return nil
}

func (c *CachingStore) DeleteEvent(ctx context.Context, userID, eventID uuid.UUID) error {
	if err := c.next.DeleteEvent(ctx, userID, eventID); err != nil {
		return err
	}
	c.invalidateUser(ctx, userID)
	return nil
}

func cacheKey(userID uuid.UUID, periodKey string) string {
	return fmt.Sprintf("events:%s:%s", userID, periodKey)
}

func (c *CachingStore) getCached(ctx context.Context, key string) ([]domain.Event, bool) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if !errors.Is(err, goredis.Nil) {
			log.Printf("events cache: get %q failed: %v", key, err)
		}
		return nil, false
	}

	var events []domain.Event
	if err := json.Unmarshal(data, &events); err != nil {
		log.Printf("events cache: unmarshal %q failed: %v", key, err)
		return nil, false
	}
	return events, true
}

func (c *CachingStore) setCached(ctx context.Context, key string, events []domain.Event) {
	data, err := json.Marshal(events)
	if err != nil {
		log.Printf("events cache: marshal %q failed: %v", key, err)
		return
	}
	if err := c.client.Set(ctx, key, data, c.ttl).Err(); err != nil {
		log.Printf("events cache: set %q failed: %v", key, err)
	}
}

func (c *CachingStore) invalidateUser(ctx context.Context, userID uuid.UUID) {
	pattern := fmt.Sprintf("events:%s:*", userID)

	var keys []string
	iter := c.client.Scan(ctx, 0, pattern, 100).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		log.Printf("events cache: scan %q failed: %v", pattern, err)
		return
	}
	if len(keys) == 0 {
		return
	}
	if err := c.client.Del(ctx, keys...).Err(); err != nil {
		log.Printf("events cache: del keys failed: %v", err)
	}
}
