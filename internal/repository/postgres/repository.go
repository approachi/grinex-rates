package postgres

import (
	"context"

	"go.opentelemetry.io/otel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/approachi/grinex-rates/internal/domain"
	"github.com/approachi/grinex-rates/internal/observability/tracing"
)

var tracer = otel.Tracer("repository")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Save(ctx context.Context, rate domain.Rate) error {
	ctx, span := tracer.Start(ctx, "repository.Save")
	defer span.End()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO rates (ask, bid, created_at) VALUES ($1, $2, $3)`,
		rate.Ask,
		rate.Bid,
		rate.CreatedAt,
	)
	if err != nil {
		tracing.RecordError(span, err)
	}
	return err
}
