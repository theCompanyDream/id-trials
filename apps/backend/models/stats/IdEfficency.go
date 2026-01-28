package stats

type IDEfficiency struct {
	TableName         string  `json:"table_name"`
	RowCount          int64   `json:"row_count"`
	AvgIDBytes        float64 `json:"avg_id_bytes"`
	TheoreticalMin    float64 `json:"theoretical_min_bytes"`
	EfficiencyPercent float64 `json:"efficiency_percent"`
	WasteFactor       float64 `json:"waste_factor"`
}
