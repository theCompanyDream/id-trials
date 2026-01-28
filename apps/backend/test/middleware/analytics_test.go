package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/theCompanyDream/id-trials/apps/backend/middleware"
	"github.com/theCompanyDream/id-trials/apps/backend/models"
	"github.com/theCompanyDream/id-trials/apps/backend/test/setup"
)

// func TestMiddlewareWithMockDB(t *testing.T) {
// 	// Setup
// 	db := setup.NewPostgresMockDB()
// 	defer setup.CleanupDB(t, db)

// 	e := echo.New()
// 	metricsMiddleware := middleware.NewMetricsMiddleware(db)
// 	e.Use(metricsMiddleware.CaptureMetrics())

// 	// Test route that uses DB
// 	e.GET("/api/v1/uuid/:id", func(c echo.Context) error {
// 		// Simulate different DB durations
// 		if c.Param("id") == "slow" {
// 			c.Set("db_duration", 100*time.Millisecond)
// 		} else {
// 			c.Set("db_duration", 10*time.Millisecond)
// 		}
// 		return c.String(http.StatusOK, "ok")
// 	})

// 	// Test 1: Normal request
// 	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/uuid/123", nil)
// 	rec1 := httptest.NewRecorder()
// 	e.ServeHTTP(rec1, req1)
// 	assert.Equal(t, http.StatusOK, rec1.Code)

// 	// Test 2: Slow request
// 	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/uuid/slow", nil)
// 	rec2 := httptest.NewRecorder()
// 	e.ServeHTTP(rec2, req2)
// 	assert.Equal(t, http.StatusOK, rec2.Code)

// 	// Wait for async saves
// 	time.Sleep(300 * time.Millisecond)

// 	// Verify both metrics were saved
// 	var metrics []models.RouteMetric
// 	db.Find(&metrics)
// 	assert.Equal(t, 2, len(metrics))

// 	// Find the slow request metric
// 	var slowMetric models.RouteMetric
// 	db.Where("db_query_duration > ?", 50).First(&slowMetric)
// 	assert.Equal(t, 100.0, slowMetric.DBQueryDuration)
// 	assert.Equal(t, "UUID", slowMetric.IDType)
// }

func TestMiddlewareIDTypeExtraction(t *testing.T) {
	db := setup.NewPostgresMockDB()
	defer setup.CleanupDB(t, db)

	e := echo.New()
	metricsMiddleware := middleware.NewMetricsMiddleware(db)
	e.Use(metricsMiddleware.CaptureMetrics())

	// Test different ID types
	routes := []struct {
		path     string
		expected string
	}{
		{"/api/v1/ulid/123", "ULID"},
		{"/api/v1/uuid/456", "UUID"},
		{"/api/v1/ksuid/789", "KSUID"},
		{"/api/v1/cuid/abc", "CUID"},
		{"/api/v1/nano/def", "NanoID"},
		{"/api/v1/snowid/ghi", "Snowflake"},
		{"/api/v1/other/xyz", "Unknown"},
	}

	for idx, route := range routes {
		t.Run(route.path, func(t *testing.T) {
			e.GET(route.path, func(c echo.Context) error {
				return c.String(http.StatusOK, "test")
			})

			req := httptest.NewRequest(http.MethodGet, route.path, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			time.Sleep(100 * time.Millisecond)

			var metric models.RouteMetric
			db.Last(&metric)
			var countAfter int64
			db.Model(&models.RouteMetric{}).Count(&countAfter)
			if route.expected == "Unknown" {
				assert.Equal(t, int64(idx), countAfter)
			} else {
				assert.Equal(t, route.expected, metric.IDType)
			}
		})
	}
}

func TestExtractIDType(t *testing.T) {
	ids := []string{"/ulidIds", "/nanoIds", "/ksuidIds", "/cuidIds", "/snowIds", "/uuid4"}
	idTypes := []string{"ULID", "NanoID", "KSUID", "CUID", "Snowflake", "UUID"}

	for idx, url := range ids {
		idType := middleware.ExtractIDType(url)
		withID := fmt.Sprintf("%s/234234", url)
		assert.Equal(t, idType, idTypes[idx])
		idType = middleware.ExtractIDType(withID)
		assert.Equal(t, idType, idTypes[idx])
	}
}
