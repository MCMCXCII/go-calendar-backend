package eventsclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (c *Client) GetEvent(ctx context.Context, id uuid.UUID) (Event, error) {
	const getEventPath = "api/v1/events"

	path := fmt.Sprintf("http://%s/%s/%s", c.host, getEventPath, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return Event{}, fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return Event{}, fmt.Errorf("client.Do: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// продолжаем ниже
	case http.StatusBadRequest:
		return Event{}, ErrValidation
	case http.StatusUnauthorized:
		return Event{}, ErrUnauthorized
	case http.StatusNotFound:
		return Event{}, ErrNotFound
	default:
		return Event{}, fmt.Errorf("request failed: status: %s", resp.Status)
	}

	var event Event
	if err := json.NewDecoder(resp.Body).Decode(&event); err != nil {
		return Event{}, fmt.Errorf("json.Decode: %w", err)
	}

	return event, nil
}
