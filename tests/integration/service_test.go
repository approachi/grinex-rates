package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/approachi/grinex-rates/internal/domain"
	"github.com/approachi/grinex-rates/internal/repository/postgres"
	"github.com/approachi/grinex-rates/internal/service"
)

type mockExchangeClient struct {
	asks []domain.OrderBookEntry
	bids []domain.OrderBookEntry
}

func (m *mockExchangeClient) FetchOrderBook(_ context.Context) ([]domain.OrderBookEntry, []domain.OrderBookEntry, error) {
	return m.asks, m.bids, nil
}

func defaultMockClient() *mockExchangeClient {
	return &mockExchangeClient{
		asks: []domain.OrderBookEntry{
			{Price: 8000},
			{Price: 8100},
			{Price: 8200},
		},
		bids: []domain.OrderBookEntry{
			{Price: 7900},
			{Price: 7800},
			{Price: 7700},
		},
	}
}

func TestService_ShouldReturnRate_WhenTopNMethod(t *testing.T) {
	pool := newTestPool(t)
	truncateRates(t, pool)

	svc := service.NewRatesService(defaultMockClient(), postgres.NewRepository(pool))

	rate, err := svc.GetRates(context.Background(), service.MethodTopN, 1, 0)
	require.NoError(t, err)

	assert.Equal(t, int64(8000), rate.Ask)
	assert.Equal(t, int64(7900), rate.Bid)
	assert.False(t, rate.CreatedAt.IsZero())

	var count int
	err = pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM rates`).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestService_ShouldReturnRate_WhenAvgNMMethod(t *testing.T) {
	pool := newTestPool(t)
	truncateRates(t, pool)

	svc := service.NewRatesService(defaultMockClient(), postgres.NewRepository(pool))

	rate, err := svc.GetRates(context.Background(), service.MethodAvgNM, 1, 3)
	require.NoError(t, err)

	assert.Equal(t, int64(8100), rate.Ask)
	assert.Equal(t, int64(7800), rate.Bid)
}

func TestService_ShouldReturnError_WhenUnknownMethod(t *testing.T) {
	pool := newTestPool(t)

	svc := service.NewRatesService(defaultMockClient(), postgres.NewRepository(pool))

	_, err := svc.GetRates(context.Background(), "unknown", 1, 0)
	assert.ErrorIs(t, err, domain.ErrUnknownMethod)
}
