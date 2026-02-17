package tenants

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrTenantNameExists = errors.New("tenant name already exists")
var ErrTenantNotFound = errors.New("tenant not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, in CreateTenantInput) (Tenant, error) {
	const q = `
		INSERT INTO tenants (name)
		VALUES ($1)
		RETURNING id, name, created_at;
	`
	var t Tenant
	err := r.pool.QueryRow(ctx, q, in.Name).Scan(&t.ID, &t.Name, &t.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return Tenant{}, ErrTenantNameExists
		}
		return Tenant{}, err
	}
	return t, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (Tenant, error) {
	const q = `
		SELECT id, name, created_at
		FROM tenants
		WHERE id = $1;
	`
	var t Tenant
	err := r.pool.QueryRow(ctx, q, id).Scan(&t.ID, &t.Name, &t.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Tenant{}, ErrTenantNotFound
		}
		return Tenant{}, err
	}
	return t, nil
}
