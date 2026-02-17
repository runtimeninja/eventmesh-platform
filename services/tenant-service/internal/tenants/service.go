package tenants

import (
	"context"
	"errors"
	"strings"
)

var ErrInvalidTenantName = errors.New("invalid tenant name")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, in CreateTenantInput) (Tenant, error) {
	in.Name = strings.TrimSpace(in.Name)
	if len(in.Name) < 2 {
		return Tenant{}, ErrInvalidTenantName
	}
	return s.repo.Create(ctx, in)
}

func (s *Service) GetByID(ctx context.Context, id string) (Tenant, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Tenant{}, ErrTenantNotFound
	}
	return s.repo.GetByID(ctx, id)
}
