package stats

import "time"

type ErrorRateTrend struct {
	TimeBucket   time.Time `json:"time_bucket"`
	RequestCount int64     `json:"request_count"`
	ErrorCount   int64     `json:"error_count"`
	ErrorRate    float64   `json:"error_rate"`
}

type PercentileTrend struct {
	TimeBucket   time.Time `json:"time_bucket"`
	RequestCount int64     `json:"request_count"`
	P50Duration  float64   `json:"p50_duration"`
	P95Duration  float64   `json:"p95_duration"`
	P99Duration  float64   `json:"p99_duration"`
}
