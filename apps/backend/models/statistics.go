package models

import "time"

type RouteStatistics struct {
	ID uint `gorm:"primaryKey"`

	// Grouping
	IDType     string    `gorm:"type:varchar(20);not null;index:idx_stats"`
	RoutePath  string    `gorm:"type:varchar(255);not null"`
	HTTPMethod string    `gorm:"type:varchar(10);not null"`
	TimeWindow time.Time `gorm:"not null;index:idx_time_window"` // Hourly/Daily bucket

	// Performance Statistics
	TotalRequests  int `gorm:"not null"`
	SuccessfulReqs int `gorm:"not null"`
	FailedReqs     int `gorm:"not null"`

	// Duration Statistics (milliseconds)
	AvgDuration float64 `gorm:"not null"`
	MinDuration float64 `gorm:"not null"`
	MaxDuration float64 `gorm:"not null"`
	P50Duration float64 `gorm:"not null"` // Median
	P95Duration float64 `gorm:"not null"` // 95th percentile
	P99Duration float64 `gorm:"not null"` // 99th percentile

	// Database Performance
	AvgDBDuration float64 `gorm:"not null"`

	// Error Rate
	ErrorRate float64 `gorm:"not null"` // Percentage

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (RouteStatistics) TableName() string {
	return "route_statistics"
}
