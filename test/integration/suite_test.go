//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"project/pkg/httpclient/authclient"
	"project/pkg/httpclient/eventclient"
)

// prepare make up
// run tests make test-integration

const testPassword = "password123"

var ctx = context.Background()

func Test_Integration(t *testing.T) {
	suite.Run(t, &Suite{})
}

type Suite struct {
	suite.Suite
	*require.Assertions

	auth   *authclient.Client
	events *eventclient.Client
}

func (s *Suite) SetupSuite() {
	s.Assertions = s.Require()

	s.waitForReady("http://localhost:8081/api/v1/auth/login")
	s.waitForReady("http://localhost:8080/api/v1/events/")

	s.auth = authclient.New(authclient.Config{Host: "localhost", Port: "8081"})
	s.events = eventclient.New(eventclient.Config{Host: "localhost", Port: "8080"})
}

func (s *Suite) waitForReady(url string) {
	client := &http.Client{Timeout: 2 * time.Second}
	deadline := time.Now().Add(30 * time.Second)

	for {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()
			return
		}
		if time.Now().After(deadline) {
			s.FailNowf("service not ready", "url=%s err=%v", url, err)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func (s *Suite) uniqueEmail() string {
	return fmt.Sprintf("integration-%d@example.com", time.Now().UnixNano())
}

func (s *Suite) registerAndLogin() *eventclient.Client {
	email := s.uniqueEmail()

	_, err := s.auth.Register(ctx, email, testPassword)
	s.NoError(err)

	token, err := s.auth.Login(ctx, email, testPassword)
	s.NoError(err)

	client := eventclient.New(eventclient.Config{Host: "localhost", Port: "8080"})
	client.SetToken(token)
	return client
}
