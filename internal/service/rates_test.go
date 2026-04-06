package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/approachi/grinex-rates/internal/domain"
)

type mockClient struct {
	asks []domain.OrderBookEntry
	bids []domain.OrderBookEntry
	err  error
}

func (m *mockClient) FetchOrderBook(_ context.Context) ([]domain.OrderBookEntry, []domain.OrderBookEntry, error) {
	return m.asks, m.bids, m.err
}

type mockRepo struct {
	saved []domain.Rate
	err   error
}

func (m *mockRepo) Save(_ context.Context, rate domain.Rate) error {
	if m.err != nil {
		return m.err
	}
	m.saved = append(m.saved, rate)
	return nil
}

var (
	testAsks = []domain.OrderBookEntry{
		{Price: 100, Volume: 10},
		{Price: 200, Volume: 20},
		{Price: 300, Volume: 30},
	}
	testBids = []domain.OrderBookEntry{
		{Price: 50, Volume: 5},
		{Price: 60, Volume: 6},
		{Price: 70, Volume: 7},
	}
)

func TestGetRates_TopN_ShouldReturnRate_WhenSuccess(t *testing.T) {
	client := &mockClient{asks: testAsks, bids: testBids}
	repo := &mockRepo{}
	svc := NewRatesService(client, repo)

	rate, err := svc.GetRates(context.Background(), MethodTopN, 1, 0)

	require.NoError(t, err)
	assert.Equal(t, int64(100), rate.Ask)
	assert.Equal(t, int64(50), rate.Bid)
	assert.False(t, rate.CreatedAt.IsZero())
	assert.Len(t, repo.saved, 1)
}

func TestGetRates_AvgNM_ShouldReturnRate_WhenSuccess(t *testing.T) {
	client := &mockClient{asks: testAsks, bids: testBids}
	repo := &mockRepo{}
	svc := NewRatesService(client, repo)

	rate, err := svc.GetRates(context.Background(), MethodAvgNM, 1, 2)

	require.NoError(t, err)
	assert.Equal(t, int64(150), rate.Ask)
	assert.Equal(t, int64(55), rate.Bid)
	assert.Len(t, repo.saved, 1)
}

func TestGetRates_ShouldReturnError_WhenClientFails(t *testing.T) {
	fetchErr := errors.New("connection refused")
	client := &mockClient{err: fetchErr}
	repo := &mockRepo{}
	svc := NewRatesService(client, repo)

	_, err := svc.GetRates(context.Background(), MethodTopN, 1, 0)

	require.ErrorIs(t, err, fetchErr)
	assert.Empty(t, repo.saved)
}

func TestGetRates_ShouldReturnError_WhenUnknownMethod(t *testing.T) {
	client := &mockClient{asks: testAsks, bids: testBids}
	repo := &mockRepo{}
	svc := NewRatesService(client, repo)

	_, err := svc.GetRates(context.Background(), "unknown", 1, 0)

	require.ErrorIs(t, err, domain.ErrUnknownMethod)
	assert.Empty(t, repo.saved)
}

func TestGetRates_ShouldReturnError_WhenRepoFails(t *testing.T) {
	client := &mockClient{asks: testAsks, bids: testBids}
	repo := &mockRepo{err: errors.New("db error")}
	svc := NewRatesService(client, repo)

	_, err := svc.GetRates(context.Background(), MethodTopN, 1, 0)

	require.Error(t, err)
}
