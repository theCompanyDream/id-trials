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
			PERCENTILE_CONT(0.25) WITHIN GROUP (ORDER BY total_duration) as quartile1,
			PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY total_duration) as median,
			PERCENTILE_CONT(0.75) WITHIN GROUP (ORDER BY total_duration) as quartile3,
            MAX(total_duration) as max_duration
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

// Get time series data (for charts)
func (r *MetricsRepository) GetErrorRateTrend(idType string) ([]stats.ErrorRateTrend, error) {
	var results []stats.ErrorRateTrend

	err := r.DB.Model(&models.RouteMetric{}).
		Select(`
			DATE_TRUNC('hour', timestamp) as time_bucket,
			COUNT(*) as request_count,
			SUM(CASE WHEN is_error THEN 1 ELSE 0 END) as error_count,
			ROUND(100.0 * SUM(CASE WHEN is_error THEN 1 ELSE 0 END) / COUNT(*), 2) as error_rate
		`).
		Where("id_type = ?", idType).
		Group("DATE_TRUNC('hour', timestamp)").
		Order("time_bucket ASC").
		Limit(20).
		Scan(&results).Error

	return results, err
}

func (r *MetricsRepository) GetIdDurationTrend(idType string) ([]stats.PercentileTrend, error) {
	var results []stats.PercentileTrend

	err := r.DB.Model(&models.RouteMetric{}).
		Select(`
			DATE_TRUNC('hour', timestamp) as time_bucket,
			COUNT(*) as request_count,
			PERCENTILE_CONT(0.50) WITHIN GROUP (ORDER BY total_duration) as p50_duration,
			PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY total_duration) as p95_duration,
			PERCENTILE_CONT(0.99) WITHIN GROUP (ORDER BY total_duration) as p99_duration
		`).
		Where("id_type = ?", idType).
		Group("DATE_TRUNC('hour', timestamp)").
		Order("time_bucket ASC").
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
