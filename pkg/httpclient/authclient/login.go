package authclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) Login(ctx context.Context, email, password string) (string, error) {
	const login_path = "api/v1/auth/login"

	path := fmt.Sprintf("http://%s/%s", c.host, login_path)

	request := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    email,
		Password: password,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("json.Marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Content-type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("client.Do: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// ниже
	case http.StatusBadRequest:
		return "", ErrValidation
	case http.StatusUnauthorized:
		return "", ErrInvalidCredentials
	default:
		return "", fmt.Errorf("request failed: status: %s", resp.Status)
	}

	response := struct {
		AccessToken string `json:"token"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("json.Decode: %w", err)
	}

	return response.AccessToken, nil
}
