package calculator

import (
	"github.com/approachi/grinex-rates/internal/domain"
)

// TopN returns the price at position n (1-based) in the order book.
func TopN(entries []domain.OrderBookEntry, n int) (int64, error) {
	if len(entries) == 0 {
		return 0, domain.ErrEmptyEntries
	}
	if n < 1 || n > len(entries) {
		return 0, domain.ErrOutOfRange
	}
	return entries[n-1].Price, nil
}

// AvgNM returns the average price between positions from and to (1-based, inclusive).
func AvgNM(entries []domain.OrderBookEntry, from, to int) (int64, error) {
	if len(entries) == 0 {
		return 0, domain.ErrEmptyEntries
	}
	if from < 1 || to > len(entries) {
		return 0, domain.ErrOutOfRange
	}
	if from > to {
		return 0, domain.ErrInvalidRange
	}

	var sum int64
	for _, e := range entries[from-1 : to] {
		sum += e.Price
	}
	return sum / int64(to-from+1), nil
}
