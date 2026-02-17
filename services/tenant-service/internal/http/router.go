package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/runtimeninja/eventmesh-platform/services/tenant-service/internal/tenants"
)

type Server struct {
	log    *slog.Logger
	tenant *tenants.Service
}

func NewServer(log *slog.Logger, tenant *tenants.Service) *Server {
	return &Server{log: log, tenant: tenant}
}

func (s *Server) Routes(readTimeout time.Duration) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Compress(5))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"ok": true})
	})

	r.Route("/v1", func(r chi.Router) {
		r.Post("/tenants", s.handleCreateTenant())
		r.Get("/tenants/{id}", s.handleGetTenant())
	})

	return r
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
