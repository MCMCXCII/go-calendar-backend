package eventclient

import (
	"net"
	"net/http"
	"time"
)

type Config struct {
	Host string `envconfig:"EVENTS_CLIENT_HOST" default:"localhost"`
	Port string `envconfig:"EVENTS_CLIENT_PORT" default:"8080"`
}

type Client struct {
	client http.Client
	host   string
	token  string
}

func New(c Config) *Client {
	return &Client{
		client: http.Client{
			Timeout: 5 * time.Second,
		},
		host: net.JoinHostPort(c.Host, c.Port),
	}
}

func (c *Client) SetToken(token string) {
	c.token = token
}
