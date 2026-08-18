package http

import (
	"context"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	v1 "project/internal/events/controller/http/v1"
	"project/internal/events/service"
	"project/pkg/token"
)

type tokenParser interface {
	ParseAccessToken(tokenString string) (token.Info, error)
}

type blackList interface {
	IsRevoked(ctx context.Context, tokenID string) (bool, error)
}

func EventsRouter(r *chi.Mux, svc *service.Service, tp tokenParser, bl blackList) {
	handler := v1.New(svc)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1/events", func(r chi.Router) {
		r.Use(v1.Auth(tp, bl))

		r.Post("/", handler.CreateEvent)
		r.Get("/", handler.ListEvents)
		r.Get("/{id}", handler.GetEvent)
		r.Put("/{id}", handler.UpdateEvent)
		r.Delete("/{id}", handler.DeleteEvent)
	})
}
