package entity

// GRCUnit represents a work unit with GRC component scores
type GRCUnit struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"` // Balai Besar, Kelas I, Kelas II, Loka, Direktorat

	// Raw scores (original values)
	SKI   float64 `json:"ski"`
	SPIP  float64 `json:"spip"`
	MRI   float64 `json:"mri"`
	IEPK  float64 `json:"iepk"`
	SPIPT float64 `json:"spipt"`
	WBKRB float64 `json:"wbkrb"`
	SAKIP float64 `json:"sakip"`

	// Percentage scores (0-100 scale) - used for radar chart
	PSKI    float64 `json:"pski"`
	PSPIP   float64 `json:"pspip"`
	PMRI    float64 `json:"pmri"`
	PIEPK   float64 `json:"piepk"`
	PWBKRB  float64 `json:"pwbkrb"`
	PSAKIP  float64 `json:"psakip"`
	Average float64 `json:"average"`
}

// ComponentScore represents a single component score with analysis
type ComponentScore struct {
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// GapAnalysis represents gap analysis for a component
type GapAnalysis struct {
	Component string  `json:"component"`
	Value     float64 `json:"value"`
	Average   float64 `json:"average"`
	Gap       float64 `json:"gap"`
	Status    string  `json:"status"` // "above", "below", "equal"
}

// WeaknessPriority represents weakness analysis for a unit
type WeaknessPriority struct {
	Component      string  `json:"component"`
	Value          float64 `json:"value"`
	GapFromAverage float64 `json:"gap_from_average"`
	Priority       string  `json:"priority"` // "high", "medium", "low"
}

// CategoryStats represents statistics for a unit category
type CategoryStats struct {
	Name    string  `json:"name"`
	Count   int     `json:"count"`
	Average float64 `json:"average"`
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
}

// PerformanceLevel represents performance distribution
type PerformanceLevel struct {
	Level string `json:"level"` // "Excellent", "Good", "Fair", "Needs Improvement"
	Count int    `json:"count"`
	Min   int    `json:"min_score"`
	Max   int    `json:"max_score"`
}

// ComponentStats represents statistics for a single component
type ComponentStats struct {
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Average      float64 `json:"average"`
	Min          float64 `json:"min"`
	Max          float64 `json:"max"`
	StdDeviation float64 `json:"std_deviation"`
	UnitsBelow80 int     `json:"units_below_80"`
}

// RadarData represents data for radar chart visualization
type RadarData struct {
	Labels        []string  `json:"labels"`
	Values        []float64 `json:"values"`
	AverageValues []float64 `json:"average_values"`
}
