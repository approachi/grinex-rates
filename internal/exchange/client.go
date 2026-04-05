package exchange

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/approachi/grinex-rates/internal/domain"
)

const (
	grinexURL = "https://grinex.io/api/v1/spot/depth"
	symbol    = "usdta7a5"
	precision = 100
)

type depthResponse struct {
	Asks []orderEntry `json:"asks"`
	Bids []orderEntry `json:"bids"`
}

type orderEntry struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
}

type GrinexClient struct {
	http *resty.Client
}

func NewGrinexClient() *GrinexClient {
	return &GrinexClient{
		http: resty.New(),
	}
}

func (c *GrinexClient) FetchOrderBook(ctx context.Context) (asks, bids []domain.OrderBookEntry, err error) {
	var resp depthResponse

	_, err = c.http.R().
		SetContext(ctx).
		SetQueryParam("symbol", symbol).
		SetResult(&resp).
		Get(grinexURL)
	if err != nil {
		return nil, nil, err
	}

	asks = make([]domain.OrderBookEntry, len(resp.Asks))
	for i, a := range resp.Asks {
		asks[i] = domain.OrderBookEntry{
			Price:  parsePrice(a.Price),
			Volume: parsePrice(a.Volume),
		}
	}

	bids = make([]domain.OrderBookEntry, len(resp.Bids))
	for i, b := range resp.Bids {
		bids[i] = domain.OrderBookEntry{
			Price:  parsePrice(b.Price),
			Volume: parsePrice(b.Volume),
		}
	}

	return asks, bids, nil
}

// parsePrice converts a decimal string to int64 in minor units (cents).
// Example: "80.85" -> 8085
func parsePrice(s string) int64 {
	parts := strings.SplitN(s, ".", 2)

	intPart, _ := strconv.ParseInt(parts[0], 10, 64)
	result := intPart * precision

	if len(parts) == 2 {
		dec := parts[1]
		for len(dec) < 2 {
			dec += "0"
		}
		dec = dec[:2]
		decVal, _ := strconv.ParseInt(dec, 10, 64)
		result += decVal
	}

	return result
}
