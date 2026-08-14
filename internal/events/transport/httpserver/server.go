package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"project/internal/platform/token"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type tokenParser interface {
	ParseAccessToken(tokenString string) (token.Info, error)
}

type blackList interface {
	IsRevoked(ctx context.Context, tokenID string) (bool, error)
}

type Server struct {
	router    *chi.Mux
	addr      string
	app       app
	validate  *validator.Validate
	ready     atomic.Bool
	token     tokenParser
	blacklist blackList
}

type Params struct {
	Addr      string
	App       app
	Token     tokenParser
	Blacklist blackList
}

func New(p Params) *Server {
	s := &Server{
		router:    chi.NewRouter(),
		addr:      p.Addr,
		app:       p.App,
		validate:  validator.New(),
		token:     p.Token,
		blacklist: p.Blacklist,
	}

	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Recoverer)

	s.router.Route("/api/v1/events", func(r chi.Router) {
		r.Use(Auth(s.token, s.blacklist))

		r.Post("/", s.handleCreateEvent)
		r.Get("/", s.handleListEvents)
		r.Get("/{id}", s.handleGetEvent)
		r.Put("/{id}", s.handleUpdateEvent)
		r.Delete("/{id}", s.handleDeleteEvent)
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

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
