package utils

import "github.com/theCompanyDream/id-trials/apps/backend/models/stats"

func CalculatePercentiles(sorted []float64) *stats.PercentileStats {
	if len(sorted) == 0 {
		return &stats.PercentileStats{}
	}

	return &stats.PercentileStats{
		P50: Percentile(sorted, 0.50),
		P75: Percentile(sorted, 0.75),
		P90: Percentile(sorted, 0.90),
		P95: Percentile(sorted, 0.95),
		P99: Percentile(sorted, 0.99),
	}
}

func Percentile(sorted []float64, p float64) float64 {
	index := int(float64(len(sorted)-1) * p)
	return sorted[index]
}
