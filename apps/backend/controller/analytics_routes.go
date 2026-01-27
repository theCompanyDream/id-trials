package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/theCompanyDream/id-trials/apps/backend/repository"
	"gorm.io/gorm"
)

type AnalyticsController struct {
	Repo *repository.MetricsRepository
}

func NewAnalyticsController(db *gorm.DB) *AnalyticsController {
	return &AnalyticsController{
		Repo: repository.NewMetricsRepository(db),
	}
}

// GetIDTypeComparison godoc
// @Summary Get performance comparison across all ID types
// @Description Returns average duration metrics for all ID types
// @Tags Analytics
// @Accept json
// @Produce json
// @Success 200 {array} stats.IDTypePerformance
// @Failure 500 {object} map[string]string
// @Router /analytics/comparison [get]
func (ac *AnalyticsController) GetIDTypeComparison(c echo.Context) error {
	results, err := ac.Repo.GetAverageDurationByIDType()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, results)
}

// GetIDTypeDetails godoc
// @Summary Get detailed performance for specific ID type
// @Description Returns performance metrics by route for a specific ID type
// @Tags Analytics
// @Accept json
// @Produce json
// @Param type path string true "ID Type" Enums(uuid, ulid, ksuid, cuid, nanoid, snowflake)
// @Success 200 {array} stats.RoutePerformance
// @Failure 500 {object} map[string]string
// @Router /analytics/details/{type} [get]
func (ac *AnalyticsController) GetIDTypeDetails(c echo.Context) error {
	idType := c.Param("type")
	results, err := ac.Repo.GetPerformanceByRoute(idType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, results)
}

// GetPercentiles godoc
// @Summary Get percentile statistics
// @Description Returns percentile statistics (p50, p95, p99) for a specific ID type
// @Tags Analytics
// @Accept json
// @Produce json
// @Param type path string true "ID Type" Enums(uuid, ulid, ksuid, cuid, nanoid, snowflake)
// @Param hours query int false "Number of hours to look back" default(24)
// @Success 200 {object} stats.PercentileStats
// @Failure 500 {object} map[string]string
// @Router /analytics/percentiles/{type} [get]
func (ac *AnalyticsController) GetPercentiles(c echo.Context) error {
	idType := c.Param("type")
	hours, _ := strconv.Atoi(c.QueryParam("hours"))
	if hours == 0 {
		hours = 24
	}

	stats, err := ac.Repo.GetPercentiles(idType, hours)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, stats)
}

// GetErrorRates godoc
// @Summary Get error rates
// @Description Returns error rates across all ID types and routes
// @Tags Analytics
// @Accept json
// @Produce json
// @Success 200 {array} stats.ErrorRate
// @Failure 500 {object} map[string]string
// @Router /analytics/errors [get]
func (ac *AnalyticsController) GetErrorRates(c echo.Context) error {
	results, err := ac.Repo.GetErrorRates()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, results)
}

// GetTimeSeries godoc
// @Summary Get time series data for charts
// @Description Returns time series performance data for a specific ID type
// @Tags Analytics
// @Accept json
// @Produce json
// @Param type path string true "ID Type" Enums(uuid, ulid, ksuid, cuid, nanoid, snowflake)
// @Param hours query int false "Number of hours to look back" default(24)
// @Success 200 {array} stats.TimeSeriesPoint
// @Failure 500 {object} map[string]string
// @Router /analytics/timeseries/{type} [get]
func (ac *AnalyticsController) GetTimeSeries(c echo.Context) error {
	idType := c.Param("type")
	hours, _ := strconv.Atoi(c.QueryParam("hours"))
	if hours == 0 {
		hours = 24
	}

	results, err := ac.Repo.GetTimeSeriesData(idType, "hour", hours)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, results)
}

// GetTableSizeData godoc
// @Summary Get table size data
// @Description Returns database table size metrics for all ID types
// @Tags Analytics
// @Accept json
// @Produce json
// @Success 200 {array} stats.TableSize
// @Failure 500 {object} map[string]string
// @Router /analytics/table-sizes [get]
func (ac *AnalyticsController) GetTableSizeData(c echo.Context) error {
	results, err := ac.Repo.GetSpecificTableSizes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, results)
}

// GetIdEfficiencyMetrics godoc
// @Summary Get ID efficiency metrics
// @Description Returns efficiency metrics comparing different ID types
// @Tags Analytics
// @Accept json
// @Produce json
// @Success 200 {array} stats.IDEfficiency
// @Failure 500 {object} map[string]string
// @Router /analytics/efficiency [get]
func (ac *AnalyticsController) GetIdEfficiencyMetrics(c echo.Context) error {
	results, err := ac.Repo.GetIdEfficiencyMetrics()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, results)
}
