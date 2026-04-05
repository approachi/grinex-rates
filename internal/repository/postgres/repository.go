package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/approachi/grinex-rates/internal/domain"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Save(ctx context.Context, rate domain.Rate) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO rates (ask, bid, created_at) VALUES ($1, $2, $3)`,
		rate.Ask,
		rate.Bid,
		rate.CreatedAt,
	)
	return err
}
