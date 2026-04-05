package calculator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/approachi/grinex-rates/internal/domain"
)

func TestTopN_ShouldReturnPriceAtPosition_WhenValidInput(t *testing.T) {
	entries := []domain.OrderBookEntry{
		{Price: 100},
		{Price: 200},
		{Price: 300},
	}

	tests := []struct {
		name      string
		n         int
		expected  int64
		expectErr error
	}{
		{"position 1", 1, 100, nil},
		{"last position", 3, 300, nil},
		{"position 0 out of range", 0, 0, domain.ErrOutOfRange},
		{"position 4 out of range", 4, 0, domain.ErrOutOfRange},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price, err := TopN(entries, tt.n)
			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, price)
		})
	}

	t.Run("empty entries", func(t *testing.T) {
		_, err := TopN([]domain.OrderBookEntry{}, 1)
		assert.ErrorIs(t, err, domain.ErrEmptyEntries)
	})
}

func TestAvgNM_ShouldReturnAveragePrice_WhenValidRange(t *testing.T) {
	entries := []domain.OrderBookEntry{
		{Price: 100},
		{Price: 200},
		{Price: 300},
		{Price: 400},
		{Price: 500},
	}

	tests := []struct {
		name      string
		from, to  int
		expected  int64
		expectErr error
	}{
		{"single position", 1, 1, 100, nil},
		{"full range", 1, 5, 300, nil},
		{"partial range", 2, 4, 300, nil},
		{"invalid range", 3, 1, 0, domain.ErrInvalidRange},
		{"out of range", 1, 6, 0, domain.ErrOutOfRange},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price, err := AvgNM(entries, tt.from, tt.to)
			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, price)
		})
	}

	t.Run("empty entries", func(t *testing.T) {
		_, err := AvgNM([]domain.OrderBookEntry{}, 1, 2)
		assert.ErrorIs(t, err, domain.ErrEmptyEntries)
	})
}
