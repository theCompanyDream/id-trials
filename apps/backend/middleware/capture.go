package middleware

import (
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/theCompanyDream/id-trials/apps/backend/models"
	"gorm.io/gorm"
)

type MetricsMiddleware struct {
	DB *gorm.DB
}

func NewMetricsMiddleware(db *gorm.DB) *MetricsMiddleware {
	return &MetricsMiddleware{DB: db}
}

func (m *MetricsMiddleware) CaptureMetrics() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Store start time in context for DB timing
			c.Set("metrics_start", start)

			// Call the handler
			err := next(c)

			// Calculate duration
			duration := time.Since(start)

			// Get DB query duration if available
			dbDuration := float64(0)
			if dbTime, ok := c.Get("db_duration").(time.Duration); ok {
				dbDuration = float64(dbTime.Milliseconds())
			}

			// Extract ID type from route
			idType := ExtractIDType(c.Path())

			// Create metric record
			metric := models.RouteMetric{
				RoutePath:       c.Path(),
				HTTPMethod:      c.Request().Method,
				IDType:          idType,
				TotalDuration:   float64(duration.Milliseconds()),
				DBQueryDuration: dbDuration,
				HandlerDuration: float64(duration.Milliseconds()) - dbDuration,
				StatusCode:      c.Response().Status,
				ResponseSize:    int(c.Response().Size),
				IsError:         err != nil || c.Response().Status >= 400,
				RequestID:       c.Response().Header().Get(echo.HeaderXRequestID),
				Timestamp:       start,
				UserAgent:       c.Request().UserAgent(),
				IPAddress:       c.RealIP(),
			}

			if err != nil {
				metric.ErrorMessage = err.Error()
			}

			// Save asynchronously to avoid slowing down response
			go m.saveMetric(metric)

			return err
		}
	}
}

func (m *MetricsMiddleware) saveMetric(metric models.RouteMetric) {
	if err := m.DB.Create(&metric).Error; err != nil {
		// Log error but don't fail the request
		println("Failed to save metric:", err.Error())
	}
}

func ExtractIDType(path string) string {
	switch {
	case strings.Contains(path, "ulid"):
		return "ULID"
	case strings.Contains(path, "uuid"):
		return "UUID"
	case strings.Contains(path, "ksuid"):
		return "KSUID"
	case strings.Contains(path, "cuid"):
		return "CUID"
	case strings.Contains(path, "nano"):
		return "NanoID"
	case strings.Contains(path, "snow"):
		return "Snowflake"
	default:
		return "Unknown"
	}
}
