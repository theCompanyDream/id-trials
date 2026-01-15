package models

import "time"

// Response structs
type IDTypePerformance struct {
	IDType       string  `json:"id_type"`
	AvgDuration  float64 `json:"avg_duration"`
	RequestCount int     `json:"request_count"`
}

type RoutePerformance struct {
	RoutePath     string  `json:"route_path"`
	HTTPMethod    string  `json:"http_method"`
	AvgDuration   float64 `json:"avg_duration"`
	MinDuration   float64 `json:"min_duration"`
	MaxDuration   float64 `json:"max_duration"`
	AvgDBDuration float64 `json:"avg_db_duration"`
	RequestCount  int     `json:"request_count"`
	ErrorCount    int     `json:"error_count"`
}

type PercentileStats struct {
	P50 float64 `json:"p50"`
	P75 float64 `json:"p75"`
	P90 float64 `json:"p90"`
	P95 float64 `json:"p95"`
	P99 float64 `json:"p99"`
}

type ErrorRate struct {
	IDType          string  `json:"id_type"`
	TotalRequests   int     `json:"total_requests"`
	ErrorCount      int     `json:"error_count"`
	ErrorPercentage float64 `json:"error_percentage"`
}

type TimeSeriesPoint struct {
	TimeBucket   time.Time `json:"time_bucket"`
	AvgDuration  float64   `json:"avg_duration"`
	RequestCount int       `json:"request_count"`
}

type TableSize struct {
	TableName string `json:"table_name"`
	Size      string `json:"size"`
}
