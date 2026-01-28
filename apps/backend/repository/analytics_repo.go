package repository

import (
	"time"

	"github.com/theCompanyDream/id-trials/apps/backend/models"
	"github.com/theCompanyDream/id-trials/apps/backend/models/stats"
	"github.com/theCompanyDream/id-trials/apps/backend/utils"
	"gorm.io/gorm"
)

type MetricsRepository struct {
	DB *gorm.DB
}

func NewMetricsRepository(db *gorm.DB) *MetricsRepository {
	return &MetricsRepository{DB: db}
}

// Get average response time by ID type
func (r *MetricsRepository) GetAverageDurationByIDType() ([]stats.IDTypePerformance, error) {
	var results []stats.IDTypePerformance

	err := r.DB.Model(&models.RouteMetric{}).
		Select("id_type, AVG(total_duration) as avg_duration, COUNT(*) as request_count").
		Where("is_error = ?", false).
		Group("id_type").
		Scan(&results).Error

	return results, err
}

// Get performance by route and operation
func (r *MetricsRepository) GetPerformanceByRoute(idType string) ([]stats.RoutePerformance, error) {
	var results []stats.RoutePerformance

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
func (r *MetricsRepository) GetPercentiles(idType string, hours int) (*map[string][]stats.PercentilePoint, error) {
	// var result []stats.PercentileStats
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	result := make(map[string][]stats.PercentilePoint)

	for _, method := range methods {
		var durations []float64
		err := r.DB.Model(&models.RouteMetric{}).
			Select("total_duration").
			Where("id_type = ? AND is_error = ? AND timestamp >= ? AND http_method = ?",
				idType, false, time.Now().Add(-time.Duration(hours)*time.Hour), method).
			Order("total_duration ASC").
			Pluck("total_duration", &durations).Error

		if err != nil {
			return nil, err
		}

		result[method] = utils.CalculatePercentiles(durations)
	}

	return &result, nil
}

// Get error rate by ID type
func (r *MetricsRepository) GetErrorRates() ([]stats.ErrorRate, error) {
	var results []stats.ErrorRate

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
func (r *MetricsRepository) GetTimeSeriesData(idType string, interval string, hours int) ([]stats.TimeSeriesPoint, error) {
	var results []stats.TimeSeriesPoint

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

func (r *MetricsRepository) GetSpecificTableSizes() ([]stats.TableSize, error) {
	var sizes []stats.TableSize

	err := r.DB.Raw(`
		SELECT
			tablename as table_name,
			REGEXP_REPLACE(
				pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)),
				'[^0-9.]', '', 'g'
			)::numeric AS size,
			pg_total_relation_size(schemaname||'.'||tablename) AS size_bytes,
			pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size_pretty
		FROM pg_tables
		WHERE schemaname = 'public'
			AND tablename LIKE 'users%'
		ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC
    `).Scan(&sizes).Error

	return sizes, err
}

func (r *MetricsRepository) GetIdEfficiencyMetrics() ([]stats.IDEfficiency, error) {
	var results []stats.IDEfficiency

	err := r.DB.Raw(`
		WITH id_stats AS (
		SELECT
			'users_uuid' AS table_name,
			COUNT(*) AS row_count,
			AVG(pg_column_size(id))::numeric AS avg_id_bytes
		FROM users_uuid

		UNION ALL

		SELECT 'users_ulid', COUNT(*), AVG(pg_column_size(id))::numeric FROM users_ulid
		UNION ALL
		SELECT 'users_cuid', COUNT(*), AVG(pg_column_size(id))::numeric FROM users_cuid
		UNION ALL
		SELECT 'users_nanoid', COUNT(*), AVG(pg_column_size(id))::numeric FROM users_nanoid
		UNION ALL
		SELECT 'users_ksuid', COUNT(*), AVG(pg_column_size(id))::numeric FROM users_ksuid
		UNION ALL
		SELECT 'users_snow', COUNT(*), AVG(pg_column_size(id))::numeric FROM users_snowflake
	)
	SELECT
		table_name,
		row_count,
		avg_id_bytes,
		-- Theoretical minimum bytes needed (log2(n) / 8)
		(LOG(2, row_count) / 8)::numeric(10,2) AS theoretical_min_bytes,
		-- Efficiency score (lower is worse, 100% is perfect)
		ROUND(((LOG(2, row_count) / 8) / avg_id_bytes * 100)::numeric, 2) AS efficiency_percent,
		-- Waste factor (how many times larger than needed)
		ROUND((avg_id_bytes / (LOG(2, row_count) / 8))::numeric, 2) AS waste_factor
		FROM id_stats
		ORDER BY efficiency_percent DESC
	`).Scan(&results).Error

	return results, err
}
