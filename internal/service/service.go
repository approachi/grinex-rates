package service

import (
	"context"

	"github.com/approachi/grinex-rates/internal/domain"
)

type RatesService interface {
	GetRates(ctx context.Context, method string, positionFrom, positionTo int) (domain.Rate, error)
}
