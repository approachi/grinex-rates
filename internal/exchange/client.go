package exchange

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

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

func NewGrinexClient(timeoutSeconds int) *GrinexClient {
	return &GrinexClient{
		http: resty.New().SetTimeout(time.Duration(timeoutSeconds) * time.Second),
	}
}

func (c *GrinexClient) FetchOrderBook(ctx context.Context) (asks, bids []domain.OrderBookEntry, err error) {
	var result depthResponse

	resp, err := c.http.R().
		SetContext(ctx).
		SetQueryParam("symbol", symbol).
		SetResult(&result).
		Get(grinexURL)
	if err != nil {
		return nil, nil, err
	}

	if resp.IsError() {
		return nil, nil, fmt.Errorf("exchange API returned status %d: %s", resp.StatusCode(), resp.String())
	}

	asks = make([]domain.OrderBookEntry, len(result.Asks))
	for i, a := range result.Asks {
		price, err := parsePrice(a.Price)
		if err != nil {
			return nil, nil, fmt.Errorf("ask[%d] price: %w", i, err)
		}
		volume, err := parsePrice(a.Volume)
		if err != nil {
			return nil, nil, fmt.Errorf("ask[%d] volume: %w", i, err)
		}
		asks[i] = domain.OrderBookEntry{Price: price, Volume: volume}
	}

	bids = make([]domain.OrderBookEntry, len(result.Bids))
	for i, b := range result.Bids {
		price, err := parsePrice(b.Price)
		if err != nil {
			return nil, nil, fmt.Errorf("bid[%d] price: %w", i, err)
		}
		volume, err := parsePrice(b.Volume)
		if err != nil {
			return nil, nil, fmt.Errorf("bid[%d] volume: %w", i, err)
		}
		bids[i] = domain.OrderBookEntry{Price: price, Volume: volume}
	}

	return asks, bids, nil
}

// parsePrice converts a decimal string to int64 in minor units (cents).
// Example: "80.85" -> 8085
func parsePrice(s string) (int64, error) {
	parts := strings.SplitN(s, ".", 2)

	intPart, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer part %q: %w", parts[0], err)
	}
	result := intPart * precision

	if len(parts) == 2 {
		dec := parts[1]
		for len(dec) < 2 {
			dec += "0"
		}
		dec = dec[:2]
		decVal, err := strconv.ParseInt(dec, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid decimal part %q: %w", dec, err)
		}
		result += decVal
	}

	return result, nil
}
