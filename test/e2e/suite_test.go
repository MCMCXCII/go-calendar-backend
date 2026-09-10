// Prepare: make up
// Run test: make test-e2e

//go:build e2e

package e2e_test

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

const testPassword = "password123"

var ctx = context.Background()

func Test_E2E(t *testing.T) {
	suite.Run(t, &Suite{})
}

type Suite struct {
	suite.Suite
	*require.Assertions
}

func (s *Suite) SetupSuite() {
	s.waitForReady("http://localhost:8081/api/v1/auth/login")
	s.waitForReady("http://localhost:8080/api/v1/events/")
}

func (s *Suite) SetupTest() {
	s.Assertions = s.Require()
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
	return fmt.Sprintf("e2e-%d@example.com", time.Now().UnixNano())
}

func (s *Suite) newAuthClient() *authclient.Client {
	return authclient.New(authclient.Config{Host: "localhost", Port: "8081"})
}

func (s *Suite) newEventsClient() *eventclient.Client {
	return eventclient.New(eventclient.Config{Host: "localhost", Port: "8080"})
}
