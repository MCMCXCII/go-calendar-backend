package eventsclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (c *Client) CreateEvent(ctx context.Context, input EventInput) (uuid.UUID, error) {
	const createEventPath = "api/v1/events/"

	path := fmt.Sprintf("http://%s/%s", c.host, createEventPath)

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
		return uuid.Nil, fmt.Errorf("json.Marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return uuid.Nil, fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return uuid.Nil, fmt.Errorf("client.Do: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusCreated:
		// продолжаем ниже
	case http.StatusBadRequest:
		return uuid.Nil, ErrValidation
	case http.StatusUnauthorized:
		return uuid.Nil, ErrUnauthorized
	default:
		return uuid.Nil, fmt.Errorf("request failed: status: %s", resp.Status)
	}

	response := struct {
		ID uuid.UUID `json:"id"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return uuid.Nil, fmt.Errorf("json.Decode: %w", err)
	}

	return response.ID, nil
}
