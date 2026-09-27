package http

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	v1 "project/internal/auth/controller/http/v1"
	"project/internal/auth/service"
	"project/pkg/logger"
	"project/pkg/metrics"
)

func AuthRouter(r *chi.Mux, svc *service.Service, tp v1.TokenParser, bl v1.BlackList, m *metrics.HTTPServer) {
	handler := v1.New(svc)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(logger.Middleware)
	r.Use(metrics.NewMiddleware(m))

	r.Handle("/metrics", promhttp.Handler())

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)

		r.Group(func(r chi.Router) {
			r.Use(v1.Auth(tp, bl))
			r.Post("/logout", handler.Logout)
		})
	})
}
