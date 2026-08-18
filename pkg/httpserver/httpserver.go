package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

type Config struct {
	Address     string        `envconfig:"HTTP_ADDRESS"      `
	Timeout     time.Duration `envconfig:"HTTP_TIMEOUT"      default:"4s"`
	IdleTimeout time.Duration `envconfig:"HTTP_IDLE_TIMEOUT" default:"60s"`
}

type Server struct {
	server *http.Server
}

func New(handler http.Handler, c Config) *Server {
	httpServer := &http.Server{
		Addr:         c.Address,
		Handler:      handler,
		ReadTimeout:  c.Timeout,
		WriteTimeout: c.Timeout,
		IdleTimeout:  c.IdleTimeout,
	}

	s := &Server{
		server: httpServer,
	}

	go s.start()

	return s
}

func (s *Server) start() {
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("http server: ListenAndServe")
	}
}

func (s *Server) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		slog.Error("http server: s.server.Shutdown")
	}

	slog.Info("http server: closed")
}
