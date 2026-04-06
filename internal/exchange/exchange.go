package exchange

import (
	"context"

	"github.com/approachi/grinex-rates/internal/domain"
)

type Client interface {
	FetchOrderBook(ctx context.Context) (asks, bids []domain.OrderBookEntry, err error)
}
