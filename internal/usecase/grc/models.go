package grc

import "sandbox/internal/domain/entity"

// OverviewResponse represents the overview dashboard response
type OverviewResponse struct {
	Statistics              OverviewStatistics        `json:"statistics"`
	TopPerformers           []UnitSummary             `json:"top_performers"`
	BottomPerformers        []UnitSummary             `json:"bottom_performers"`
	WeakestComponents       []entity.ComponentStats   `json:"weakest_components"`
	CategorySummary         []entity.CategoryStats    `json:"category_summary"`
	PerformanceDistribution []entity.PerformanceLevel `json:"performance_distribution"`
	LastUpdated             string                    `json:"last_updated"`
}

// OverviewStatistics represents basic statistics
type OverviewStatistics struct {
	TotalUnits   int     `json:"total_units"`
	AverageScore float64 `json:"average_score"`
	Median       float64 `json:"median"`
	StdDeviation float64 `json:"std_deviation"`
}

// UnitSummary represents a summary of a unit
type UnitSummary struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Average  float64 `json:"average"`
	Rank     int     `json:"rank"`
}

// UnitListResponse represents the list units response
type UnitListResponse struct {
	Units      []UnitListItem `json:"units"`
	TotalCount int            `json:"total_count"`
}

// UnitListItem represents a unit in the list
type UnitListItem struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Category   string     `json:"category"`
	Average    float64    `json:"average"`
	Rank       int        `json:"rank"`
	Percentile float64    `json:"percentile"`
	Scores     UnitScores `json:"scores"`
}

// UnitDetailResponse represents detailed unit information
type UnitDetailResponse struct {
	Unit               UnitInfo                `json:"unit"`
	Scores             UnitScores              `json:"scores"`
	RadarData          entity.RadarData        `json:"radar_data"`
	GapAnalysis        []entity.GapAnalysis    `json:"gap_analysis"`
	Weakness           entity.WeaknessPriority `json:"weakness"`
	Strength           entity.WeaknessPriority `json:"strength"`
	CategoryComparison *CategoryComparison     `json:"category_comparison,omitempty"`
}

// UnitInfo represents basic unit information
type UnitInfo struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	Average    float64 `json:"average"`
	Rank       int     `json:"rank"`
	Percentile float64 `json:"percentile"`
}

// UnitScores represents all component scores
type UnitScores struct {
	PSKI   float64 `json:"pski"`
	PSPIP  float64 `json:"pspip"`
	PMRI   float64 `json:"pmri"`
	PIEPK  float64 `json:"piepk"`
	PWBKRB float64 `json:"pwbkrb"`
	PSAKIP float64 `json:"psakip"`
}

// CategoryComparison represents comparison with category average
type CategoryComparison struct {
	CategoryName    string  `json:"category_name"`
	CategoryAverage float64 `json:"category_average"`
	GapToCategory   float64 `json:"gap_to_category"`
}

// CompareResponse represents the comparison response
type CompareResponse struct {
	Units       []CompareUnit    `json:"units"`
	RadarLabels []string         `json:"radar_labels"`
	Average     entity.RadarData `json:"average"`
}

// CompareUnit represents a unit in comparison
type CompareUnit struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Category string    `json:"category"`
	Average  float64   `json:"average"`
	Rank     int       `json:"rank"`
	Values   []float64 `json:"values"`
}

// CategoryResponse represents category breakdown response
type CategoryResponse struct {
	Categories []CategoryDetail `json:"categories"`
}

// CategoryDetail represents detailed category information
type CategoryDetail struct {
	Name       string       `json:"name"`
	Count      int          `json:"count"`
	Average    float64      `json:"average"`
	Min        float64      `json:"min"`
	Max        float64      `json:"max"`
	TopUnit    *UnitSummary `json:"top_unit,omitempty"`
	BottomUnit *UnitSummary `json:"bottom_unit,omitempty"`
}
