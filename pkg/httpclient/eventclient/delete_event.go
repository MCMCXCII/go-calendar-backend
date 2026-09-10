package eventclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (c *Client) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	const deleteEventPath = "api/v1/events"

	path := fmt.Sprintf("http://%s/%s/%s", c.host, deleteEventPath, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, path, http.NoBody)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return ErrValidation
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusNotFound:
		return ErrNotFound
	default:
		return fmt.Errorf("request failed: status: %s", resp.Status)
	}
}
