package stats

import "time"

// Response structs
type IDTypePerformance struct {
	IDType       string  `json:"id_type"`
	AvgDuration  float64 `json:"avg_duration"`
	RequestCount int     `json:"request_count"`
}

type RoutePerformance struct {
	RoutePath   string  `json:"route_path"`
	HTTPMethod  string  `json:"http_method"`
	AvgDuration float64 `json:"avg_duration"`
	MinDuration float64 `json:"min_duration"`
	MaxDuration float64 `json:"max_duration"`
	Quartile1   float64 `json:"quartile1"`
	Median      float64 `json:"median"`
	Quartile3   float64 `json:"quartile3"`
}

type PercentilePoint struct {
	Percentile string  `json:"percentile"`
	Value      float64 `json:"value"`
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
	TableName  string `json:"table_name"`
	Size       int64  `json:"size"`        // Numeric size for charts
	SizeBytes  int64  `json:"size_bytes"`  // Numeric bytes for charts
	SizePretty string `json:"size_pretty"` // Human-readable (from PostgreSQL)
}
