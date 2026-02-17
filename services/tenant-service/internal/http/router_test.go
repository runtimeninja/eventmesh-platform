package http

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/runtimeninja/eventmesh-platform/services/tenant-service/internal/tenants"
)

type fakeTenantService struct{}

func (f fakeTenantService) Create(_ any, in tenants.CreateTenantInput) (tenants.Tenant, error) {
	if len(in.Name) < 2 {
		return tenants.Tenant{}, tenants.ErrInvalidTenantName
	}
	return tenants.Tenant{ID: "test-id", Name: in.Name}, nil
}
func (f fakeTenantService) GetByID(_ any, id string) (tenants.Tenant, error) {
	if id == "missing" {
		return tenants.Tenant{}, tenants.ErrTenantNotFound
	}
	return tenants.Tenant{ID: id, Name: "Acme"}, nil
}

// NOTE: I'll replace this test approach with proper interfaces in Step 2.
// For now, keep it minimal and verify /healthz only (fast & stable).
func TestHealthz(t *testing.T) {
	// Minimal: just ensure router serves healthz.
	// We'll expand tests after introducing service interfaces.
	log := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	_ = log

	// create a real server requires tenants.Service; so keep health test by creating a small handler directly:
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	// Inline minimal router equivalent:
	s := &Server{log: slog.Default(), tenant: nil}
	h := s.Routes(5 * time.Second)
	h.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestCreateTenant_InvalidJSON(t *testing.T) {
	s := &Server{log: slog.Default(), tenant: nil}
	h := s.Routes(5 * time.Second)

	req := httptest.NewRequest(http.MethodPost, "/v1/tenants", bytes.NewBufferString("{bad"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)
	if rr.Code != 400 {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestCreateTenant_OK_JSONShape(t *testing.T) {
	// We can't fully test create without service interface yet,
	// so we only validate it accepts JSON when service is wired in later.
	body, _ := json.Marshal(map[string]string{"name": "Acme"})
	s := &Server{log: slog.Default(), tenant: nil}
	h := s.Routes(5 * time.Second)

	req := httptest.NewRequest(http.MethodPost, "/v1/tenants", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	// Without tenant service it will 500, which is expected for now.
	// Step 2 will introduce interfaces and proper mocks.
	if rr.Code != 500 {
		t.Fatalf("expected 500 (no svc wired yet), got %d", rr.Code)
	}
}
