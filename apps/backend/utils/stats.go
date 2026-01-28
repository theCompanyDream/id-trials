package utils

import "github.com/theCompanyDream/id-trials/apps/backend/models/stats"

func CalculatePercentiles(sorted []float64) []stats.PercentilePoint {
	if len(sorted) == 0 {
		return []stats.PercentilePoint{}
	}

	return []stats.PercentilePoint{
		{Percentile: "P50", Value: Percentile(sorted, 0.50)},
		{Percentile: "P75", Value: Percentile(sorted, 0.75)},
		{Percentile: "P90", Value: Percentile(sorted, 0.90)},
		{Percentile: "P95", Value: Percentile(sorted, 0.95)},
		{Percentile: "P99", Value: Percentile(sorted, 0.99)},
	}
}

func Percentile(sorted []float64, p float64) float64 {
	index := int(float64(len(sorted)-1) * p)
	return sorted[index]
}
