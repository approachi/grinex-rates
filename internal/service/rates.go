package service

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"

	"github.com/approachi/grinex-rates/internal/calculator"
	"github.com/approachi/grinex-rates/internal/domain"
	"github.com/approachi/grinex-rates/internal/exchange"
	"github.com/approachi/grinex-rates/internal/observability/metrics"
	"github.com/approachi/grinex-rates/internal/observability/tracing"
	"github.com/approachi/grinex-rates/internal/repository"
)

// methods to use in calculate.
const (
	MethodTopN  = "topN"
	MethodAvgNM = "avgNM"
)

var tracer = otel.Tracer("service")

type ratesService struct {
	client exchange.Client
	repo   repository.RatesRepository
}

func NewRatesService(client exchange.Client, repo repository.RatesRepository) RatesService {
	return &ratesService{
		client: client,
		repo:   repo,
	}
}

func (s *ratesService) GetRates(ctx context.Context, method string, positionFrom, positionTo int) (domain.Rate, error) {
	ctx, span := tracer.Start(ctx, "service.GetRates")
	defer span.End()

	asks, bids, err := s.fetchOrderBook(ctx)
	if err != nil {
		tracing.RecordError(span, err)
		return domain.Rate{}, err
	}

	ask, err := calculate(method, asks, positionFrom, positionTo)
	if err != nil {
		tracing.RecordError(span, err)
		return domain.Rate{}, err
	}

	bid, err := calculate(method, bids, positionFrom, positionTo)
	if err != nil {
		tracing.RecordError(span, err)
		return domain.Rate{}, err
	}

	rate := domain.Rate{
		Ask:       ask,
		Bid:       bid,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repo.Save(ctx, rate); err != nil {
		tracing.RecordError(span, err)
		return domain.Rate{}, err
	}

	metrics.RatesAsk.Set(float64(rate.Ask))
	metrics.RatesBid.Set(float64(rate.Bid))

	return rate, nil
}

func (s *ratesService) fetchOrderBook(ctx context.Context) (asks, bids []domain.OrderBookEntry, err error) {
	ctx, span := tracer.Start(ctx, "exchange.FetchOrderBook")
	defer span.End()

	asks, bids, err = s.client.FetchOrderBook(ctx)
	if err != nil {
		tracing.RecordError(span, err)
	}
	return
}

func calculate(method string, entries []domain.OrderBookEntry, positionFrom, positionTo int) (int64, error) {
	switch method {
	case MethodTopN:
		return calculator.TopN(entries, positionFrom)
	case MethodAvgNM:
		return calculator.AvgNM(entries, positionFrom, positionTo)
	default:
		return 0, domain.ErrUnknownMethod
	}
}
