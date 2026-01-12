// repository/metrics_test.go (same package, no _test suffix)
package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculatePercentiles(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected PercentileStats
	}{
		{
			name:     "empty slice",
			input:    []float64{},
			expected: PercentileStats{},
		},
		{
			name:  "single value",
			input: []float64{100},
			expected: PercentileStats{
				P50: 100,
				P75: 100,
				P90: 100,
				P95: 100,
				P99: 100,
			},
		},
		// Add more test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculatePercentiles(tt.input)
			assert.Equal(t, &tt.expected, result)
		})
	}
}
