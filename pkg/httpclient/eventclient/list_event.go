package eventsclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) ListEvents(ctx context.Context, filter ListFilter) ([]Event, error) {
	const listEventsPath = "api/v1/events/"

	query := url.Values{}
	switch {
	case filter.Day != "":
		query.Set("day", filter.Day)
	case filter.Week != "":
		query.Set("week", filter.Week)
	case filter.Month != "":
		query.Set("month", filter.Month)
	case !filter.From.IsZero() && !filter.To.IsZero():
		query.Set("from", filter.From.Format(timeLayout))
		query.Set("to", filter.To.Format(timeLayout))
	}

	path := fmt.Sprintf("http://%s/%s?%s", c.host, listEventsPath, query.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// продолжаем ниже
	case http.StatusBadRequest:
		return nil, ErrValidation
	case http.StatusUnauthorized:
		return nil, ErrUnauthorized
	default:
		return nil, fmt.Errorf("request failed: status: %s", resp.Status)
	}

	response := struct {
		Events []Event `json:"events"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("json.Decode: %w", err)
	}

	return response.Events, nil
}
