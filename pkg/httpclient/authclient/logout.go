package authclient

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) Logout(ctx context.Context) error {
	const logoutPath = "api/v1/auth/logout"

	path := fmt.Sprintf("http://%s/%s", c.host, logoutPath)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, http.NoBody)
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
	case http.StatusUnauthorized:
		return ErrUnauthorized
	default:
		return fmt.Errorf("request failed: status: %s", resp.Status)
	}
}
