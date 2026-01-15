package models

import (
	"time"
)

type RouteMetric struct {
	ID uint `gorm:"primaryKey"`

	// Route Information
	RoutePath  string `gorm:"type:varchar(255);not null;index:idx_route_metrics"`
	HTTPMethod string `gorm:"type:varchar(10);not null"`
	IDType     string `gorm:"type:varchar(20);not null;index:idx_id_type"` // ULID, UUID, KSUID, etc.

	// Timing Metrics (in milliseconds)
	TotalDuration   float64 `gorm:"not null"` // Total request time
	DBQueryDuration float64 `gorm:"not null"` // Database query time only
	HandlerDuration float64 `gorm:"not null"` // Handler processing time

	// Response Information
	StatusCode   int `gorm:"not null;index:idx_status"`
	ResponseSize int `gorm:"default:0"` // Response body size in bytes

	// Error Tracking
	IsError      bool   `gorm:"default:false;index:idx_error"`
	ErrorMessage string `gorm:"type:text"`

	// Request Context
	RequestID string    `gorm:"type:varchar(100);index:idx_request_id"`
	Timestamp time.Time `gorm:"not null;index:idx_timestamp"`

	// Optional metadata
	UserAgent string `gorm:"type:varchar(255)"`
	IPAddress string `gorm:"type:varchar(45)"`
}

func (RouteMetric) TableName() string {
	return "route_metrics"
}
