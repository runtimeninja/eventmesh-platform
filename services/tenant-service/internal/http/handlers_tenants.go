package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/runtimeninja/eventmesh-platform/services/tenant-service/internal/tenants"
)

type errResp struct {
	Error string `json:"error"`
}

func (s *Server) handleCreateTenant() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var in tenants.CreateTenantInput
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeJSON(w, 400, errResp{Error: "invalid json"})
			return
		}

		t, err := s.tenant.Create(r.Context(), in)
		if err != nil {
			switch {
			case errors.Is(err, tenants.ErrInvalidTenantName):
				writeJSON(w, 400, errResp{Error: err.Error()})
			case errors.Is(err, tenants.ErrTenantNameExists):
				writeJSON(w, 409, errResp{Error: "tenant name already exists"})
			default:
				s.log.Error("create tenant failed", slog.Any("err", err))
				writeJSON(w, 500, errResp{Error: "internal error"})
			}
			return
		}

		s.log.Info("tenant created", slog.String("tenant_id", t.ID), slog.Duration("took", time.Since(start)))
		writeJSON(w, 201, t)
	}
}

func (s *Server) handleGetTenant() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		t, err := s.tenant.GetByID(r.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, tenants.ErrTenantNotFound):
				writeJSON(w, 404, errResp{Error: "not found"})
			default:
				s.log.Error("get tenant failed", slog.Any("err", err))
				writeJSON(w, 500, errResp{Error: "internal error"})
			}
			return
		}

		writeJSON(w, 200, t)
	}
}
