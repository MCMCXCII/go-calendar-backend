package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"project/internal/auth/service"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type app interface {
	Register(ctx context.Context, req service.RegisterParams) (service.RegisterResult, error)
	Login(ctx context.Context, req service.LoginParams) (service.LoginResult, error)
}

type Server struct {
	router   *chi.Mux
	addr     string
	app      app
	validate *validator.Validate
	ready    atomic.Bool
}

type Params struct {
	Addr string
	App  app
}

func New(p Params) *Server {
	s := &Server{router: chi.NewRouter(), addr: p.Addr, app: p.App, validate: validator.New()}

	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Recoverer)

	s.router.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", s.handleRegister)
		r.Post("/login", s.handleLogin)
	})
	return s
}

func (s *Server) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:    s.addr,
		Handler: s.router,
	}

	s.ready.Store(true)
	defer s.ready.Store(false)

	shutdownDone := make(chan struct{})

	go func() {
		defer close(shutdownDone)

		<-ctx.Done()
		s.ready.Store(false)

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
		}
	}()

	err := server.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		<-shutdownDone
		return nil
	}

	if err != nil {
		return fmt.Errorf("serve HTTP: %w", err)
	}

	return nil
}
