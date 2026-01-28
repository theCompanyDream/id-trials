package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/theCompanyDream/id-trials/apps/backend/middleware"
	"github.com/theCompanyDream/id-trials/apps/backend/models"
	"github.com/theCompanyDream/id-trials/apps/backend/test/setup"
	"gorm.io/gorm"
)

func TestCaptureMetrics(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/ulid/123", nil)
	req.Header.Set("User-Agent", "test-agent")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock DB
	mockDB := setup.NewPostgresMockDB()
	middleware := middleware.NewMetricsMiddleware(mockDB)

	// Create a test handler
	testHandler := func(c echo.Context) error {
		// Simulate DB operation
		c.Set("db_duration", 50*time.Millisecond)
		return c.String(http.StatusOK, "test response")
	}

	// Wrap handler with middleware
	middlewareFunc := middleware.CaptureMetrics()
	handler := middlewareFunc(testHandler)

	// Expect DB.Create to be called with specific metric
	mockDB.On("Create", mock.MatchedBy(func(metric interface{}) bool {
		m, ok := metric.(models.RouteMetric)
		if !ok {
			return false
		}
		return m.RoutePath == "/api/v1/ulid/123" &&
			m.HTTPMethod == "GET" &&
			m.IDType == "ULID" &&
			m.UserAgent == "test-agent" &&
			m.IsError == false
	})).Return(&gorm.DB{})

	// Execute
	err := handler(c)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Give async goroutine time to execute
	time.Sleep(100 * time.Millisecond)

	mockDB.AssertExpectations(t)
}

func TestExtractIDType(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/api/v1/ulid/123", "ULID"},
		{"/api/v1/uuid/123", "UUID"},
		{"/api/v1/ksuid/123", "KSUID"},
		{"/api/v1/cuid/123", "CUID"},
		{"/api/v1/nano/123", "NanoID"},
		{"/api/v1/snow/123", "Snowflake"},
		{"/api/v1/other/123", "Unknown"},
		{"/ulid/test", "ULID"},
		{"ulid/test", "ULID"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := extractIDType(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}
