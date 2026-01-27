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

// Get performance comparison across all ID types
func (ac *AnalyticsController) GetIDTypeComparison(c echo.Context) error {
	results, err := ac.Repo.GetAverageDurationByIDType()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, results)
}

// Get detailed performance for specific ID type
func (ac *AnalyticsController) GetIDTypeDetails(c echo.Context) error {
	idType := c.Param("type")
	results, err := ac.Repo.GetPerformanceByRoute(idType)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, results)
}

// Get percentile statistics
func (ac *AnalyticsController) GetPercentiles(c echo.Context) error {
	idType := c.Param("type")
	hours, _ := strconv.Atoi(c.QueryParam("hours"))
	if hours == 0 {
		hours = 24
	}

	stats, err := ac.Repo.GetPercentiles(idType, hours)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, stats)
}

// Get error rates
func (ac *AnalyticsController) GetErrorRates(c echo.Context) error {
	results, err := ac.Repo.GetErrorRates()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, results)
}

// Get time series data for charts
func (ac *AnalyticsController) GetTimeSeries(c echo.Context) error {
	idType := c.Param("type")
	hours, _ := strconv.Atoi(c.QueryParam("hours"))
	if hours == 0 {
		hours = 24
	}

	results, err := ac.Repo.GetTimeSeriesData(idType, "hour", hours)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, results)
}

func (ac *AnalyticsController) GetTableSizeData(c echo.Context) error {
	results, err := ac.Repo.GetSpecificTableSizes()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, results)
}

func (ac *AnalyticsController) GetIdEfficiencyMetrics(c echo.Context) error {
	results, err := ac.Repo.GetIdEfficiencyMetrics()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, results)
}
