package eventsclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (c *Client) UpdateEvent(ctx context.Context, id uuid.UUID, input EventInput) error {
	const updateEventPath = "api/v1/events"

	path := fmt.Sprintf("http://%s/%s/%s", c.host, updateEventPath, id)

	request := struct {
		Title       string `json:"title"`
		Type        string `json:"type"`
		CustomType  string `json:"custom_type,omitempty"`
		Description string `json:"description,omitempty"`
		StartTime   string `json:"start_time"`
		EndTime     string `json:"end_time"`
	}{
		Title:       input.Title,
		Type:        input.Type,
		CustomType:  input.CustomType,
		Description: input.Description,
		StartTime:   input.StartTime.Format(timeLayout),
		EndTime:     input.EndTime.Format(timeLayout),
	}

	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
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
