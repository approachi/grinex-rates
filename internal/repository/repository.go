package repository

import (
	"context"

	"github.com/approachi/grinex-rates/internal/domain"
)

type RatesRepository interface {
	Save(ctx context.Context, rate domain.Rate) error
}
