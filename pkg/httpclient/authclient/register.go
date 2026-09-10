package authclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (c *Client) Register(ctx context.Context, email, password string) (uuid.UUID, error) {
	const register_path = "api/v1/auth/register"

	path := fmt.Sprintf("http://%s/%s", c.host, register_path)

	request := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    email,
		Password: password,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return uuid.Nil, fmt.Errorf("json.Marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return uuid.Nil, fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Content-type", "application/json")

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
	case http.StatusConflict:
		return uuid.Nil, ErrEmailAlreadyExists
	default:
		return uuid.Nil, fmt.Errorf("request failed: status: %s", resp.Status)
	}

	response := struct {
		UserID uuid.UUID `json:"user_id"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return uuid.Nil, fmt.Errorf("json.Decode: %w", err)
	}

	return response.UserID, nil
}
