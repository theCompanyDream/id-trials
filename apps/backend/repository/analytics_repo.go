package repository

import (
	"time"

	"github.com/theCompanyDream/id-trials/apps/backend/models"
	"gorm.io/gorm"
)

type MetricsRepository struct {
	DB *gorm.DB
}

func NewMetricsRepository(db *gorm.DB) *MetricsRepository {
	return &MetricsRepository{DB: db}
}

// Get average response time by ID type
func (r *MetricsRepository) GetAverageDurationByIDType() ([]models.IDTypePerformance, error) {
	var results []models.IDTypePerformance

	err := r.DB.Model(&models.RouteMetric{}).
		Select("id_type, AVG(total_duration) as avg_duration, COUNT(*) as request_count").
		Where("is_error = ?", false).
		Group("id_type").
		Scan(&results).Error

	return results, err
}

// Get performance by route and operation
func (r *MetricsRepository) GetPerformanceByRoute(idType string) ([]models.RoutePerformance, error) {
	var results []models.RoutePerformance

	err := r.DB.Model(&models.RouteMetric{}).
		Select(`
            route_path,
            http_method,
            AVG(total_duration) as avg_duration,
            MIN(total_duration) as min_duration,
            MAX(total_duration) as max_duration,
            AVG(db_query_duration) as avg_db_duration,
            COUNT(*) as request_count,
            SUM(CASE WHEN is_error THEN 1 ELSE 0 END) as error_count
        `).
		Where("id_type = ?", idType).
		Group("route_path, http_method").
		Scan(&results).Error

	return results, err
}

// Get percentile performance
func (r *MetricsRepository) GetPercentiles(idType string, hours int) (*models.PercentileStats, error) {
	var durations []float64

	err := r.DB.Model(&models.RouteMetric{}).
		Select("total_duration").
		Where("id_type = ? AND is_error = ? AND timestamp >= ?",
			idType, false, time.Now().Add(-time.Duration(hours)*time.Hour)).
		Order("total_duration ASC").
		Pluck("total_duration", &durations).Error

	if err != nil {
		return nil, err
	}

	return calculatePercentiles(durations), nil
}

// Get error rate by ID type
func (r *MetricsRepository) GetErrorRates() ([]models.ErrorRate, error) {
	var results []models.ErrorRate

	err := r.DB.Model(&models.RouteMetric{}).
		Select(`
            id_type,
            COUNT(*) as total_requests,
            SUM(CASE WHEN is_error THEN 1 ELSE 0 END) as error_count,
            (SUM(CASE WHEN is_error THEN 1 ELSE 0 END) * 100.0 / COUNT(*)) as error_percentage
        `).
		Group("id_type").
		Scan(&results).Error

	return results, err
}

// Get time series data (for charts)
func (r *MetricsRepository) GetTimeSeriesData(idType string, interval string, hours int) ([]models.TimeSeriesPoint, error) {
	var results []models.TimeSeriesPoint

	// PostgreSQL-specific for hourly grouping
	query := `
        SELECT
            DATE_TRUNC('hour', timestamp) as time_bucket,
            AVG(total_duration) as avg_duration,
            COUNT(*) as request_count
        FROM route_metrics
        WHERE id_type = ?
        AND timestamp >= ?
        GROUP BY time_bucket
        ORDER BY time_bucket ASC
    `

	err := r.DB.Raw(query, idType, time.Now().Add(-time.Duration(hours)*time.Hour)).
		Scan(&results).Error

	return results, err
}

func calculatePercentiles(sorted []float64) *models.PercentileStats {
	if len(sorted) == 0 {
		return &models.PercentileStats{}
	}

	return &models.PercentileStats{
		P50: percentile(sorted, 0.50),
		P75: percentile(sorted, 0.75),
		P90: percentile(sorted, 0.90),
		P95: percentile(sorted, 0.95),
		P99: percentile(sorted, 0.99),
	}
}

func percentile(sorted []float64, p float64) float64 {
	index := int(float64(len(sorted)-1) * p)
	return sorted[index]
}

func (r *MetricsRepository) GetSpecificTableSizes() ([]models.TableSize, error) {
	var sizes []models.TableSize

	err := r.DB.Raw(`
		SELECT
			tablename as table_name,
			pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
		FROM pg_tables
		WHERE schemaname = 'public'
			AND tablename like 'users%'
		ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC
    `).Scan(&sizes).Error

	return sizes, err
}
