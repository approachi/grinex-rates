package exchange

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePrice_ShouldConvertDecimalToMinorUnits_WhenDecimalString(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"80.85", 8085},
		{"80.00", 8000},
		{"80.9", 8090},
		{"80", 8000},
		{"0.01", 1},
		{"100.99", 10099},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parsePrice(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestParsePrice_ShouldReturnError_WhenInvalidString(t *testing.T) {
	tests := []string{"abc", "12.ab", ""}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := parsePrice(input)
			assert.Error(t, err)
		})
	}
}
