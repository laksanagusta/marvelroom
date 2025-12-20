package grc

import (
	"math"
	"sort"

	"sandbox/internal/domain/entity"
)

// Repository provides access to GRC data and analytics
type Repository struct {
	provider *DataProvider
}

// NewRepository creates a new GRC repository
func NewRepository(provider *DataProvider) *Repository {
	return &Repository{provider: provider}
}

// FetchData fetches units and average from spreadsheet (call once at usecase level)
func (r *Repository) FetchData() (units []entity.GRCUnit, avg entity.GRCUnit, err error) {
	// Use FetchAll to fetch both units and average in a single call
	return r.provider.FetchAll()
}

// GetUnitByID returns a unit by its ID from provided units slice
func (r *Repository) GetUnitByID(units []entity.GRCUnit, id int) *entity.GRCUnit {
	for _, unit := range units {
		if unit.ID == id {
			return &unit
		}
	}
	return nil
}

// GetUnitsByCategory returns units filtered by category
func (r *Repository) GetUnitsByCategory(units []entity.GRCUnit, category string) []entity.GRCUnit {
	var result []entity.GRCUnit
	for _, unit := range units {
		if unit.Category == category {
			result = append(result, unit)
		}
	}
	return result
}

// GetTopPerformers returns top N performers by average score
func (r *Repository) GetTopPerformers(units []entity.GRCUnit, n int) []entity.GRCUnit {
	sorted := make([]entity.GRCUnit, len(units))
	copy(sorted, units)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Average > sorted[j].Average
	})
	if n > len(sorted) {
		n = len(sorted)
	}
	return sorted[:n]
}

// GetBottomPerformers returns bottom N performers by average score
func (r *Repository) GetBottomPerformers(units []entity.GRCUnit, n int) []entity.GRCUnit {
	sorted := make([]entity.GRCUnit, len(units))
	copy(sorted, units)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Average < sorted[j].Average
	})
	if n > len(sorted) {
		n = len(sorted)
	}
	return sorted[:n]
}

// GetUnitRank returns the rank of a unit (1 = best)
func (r *Repository) GetUnitRank(units []entity.GRCUnit, id int) int {
	sorted := make([]entity.GRCUnit, len(units))
	copy(sorted, units)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Average > sorted[j].Average
	})
	for i, unit := range sorted {
		if unit.ID == id {
			return i + 1
		}
	}
	return -1
}

// GetPercentile returns the percentile rank of a unit
func (r *Repository) GetPercentile(units []entity.GRCUnit, id int) float64 {
	rank := r.GetUnitRank(units, id)
	if rank == -1 {
		return 0
	}
	return math.Round((float64(len(units)-rank)/float64(len(units)))*10000) / 100
}

// CalculateStatistics calculates overall statistics
func (r *Repository) CalculateStatistics(units []entity.GRCUnit) (avg, median, stdDev float64) {
	n := len(units)
	if n == 0 {
		return 0, 0, 0
	}

	// Calculate average
	var sum float64
	for _, unit := range units {
		sum += unit.Average
	}
	avg = math.Round(sum/float64(n)*100) / 100

	// Calculate median
	sorted := make([]float64, n)
	for i, unit := range units {
		sorted[i] = unit.Average
	}
	sort.Float64s(sorted)
	if n%2 == 0 {
		median = (sorted[n/2-1] + sorted[n/2]) / 2
	} else {
		median = sorted[n/2]
	}
	median = math.Round(median*100) / 100

	// Calculate standard deviation
	var sumSquares float64
	for _, val := range sorted {
		diff := val - avg
		sumSquares += diff * diff
	}
	stdDev = math.Round(math.Sqrt(sumSquares/float64(n))*100) / 100

	return avg, median, stdDev
}

// GetComponentStats returns statistics for each component
func (r *Repository) GetComponentStats(units []entity.GRCUnit) []entity.ComponentStats {
	if len(units) == 0 {
		return nil
	}

	components := []struct {
		code string
		name string
		get  func(u entity.GRCUnit) float64
	}{
		{"PSKI", "Penilaian SKI", func(u entity.GRCUnit) float64 { return u.PSKI }},
		{"PSPIP", "Penilaian SPIP", func(u entity.GRCUnit) float64 { return u.PSPIP }},
		{"PMRI", "Penilaian MRI", func(u entity.GRCUnit) float64 { return u.PMRI }},
		{"PIEPK", "Penilaian IEPK", func(u entity.GRCUnit) float64 { return u.PIEPK }},
		{"PWBKRB", "Penilaian WBK/RB", func(u entity.GRCUnit) float64 { return u.PWBKRB }},
		{"PSAKIP", "Penilaian SAKIP", func(u entity.GRCUnit) float64 { return u.PSAKIP }},
	}

	var result []entity.ComponentStats
	for _, comp := range components {
		var values []float64
		for _, unit := range units {
			values = append(values, comp.get(unit))
		}

		var sum float64
		min := values[0]
		max := values[0]
		for _, v := range values {
			sum += v
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
		}
		avg := sum / float64(len(values))

		var sumSquares float64
		for _, v := range values {
			diff := v - avg
			sumSquares += diff * diff
		}
		stdDev := math.Sqrt(sumSquares / float64(len(values)))

		below80 := 0
		for _, v := range values {
			if v < 80 {
				below80++
			}
		}

		result = append(result, entity.ComponentStats{
			Code:         comp.code,
			Name:         comp.name,
			Average:      math.Round(avg*100) / 100,
			Min:          math.Round(min*100) / 100,
			Max:          math.Round(max*100) / 100,
			StdDeviation: math.Round(stdDev*100) / 100,
			UnitsBelow80: below80,
		})
	}

	return result
}

// GetCategoryStats returns statistics per category
func (r *Repository) GetCategoryStats(units []entity.GRCUnit) []entity.CategoryStats {
	categories := []string{"Balai Besar", "Kelas I", "Kelas II", "Loka", "Direktorat"}
	var result []entity.CategoryStats

	for _, cat := range categories {
		catUnits := r.GetUnitsByCategory(units, cat)
		if len(catUnits) == 0 {
			continue
		}

		var sum, min, max float64
		min = catUnits[0].Average
		max = catUnits[0].Average
		for _, u := range catUnits {
			sum += u.Average
			if u.Average < min {
				min = u.Average
			}
			if u.Average > max {
				max = u.Average
			}
		}

		result = append(result, entity.CategoryStats{
			Name:    cat,
			Count:   len(catUnits),
			Average: math.Round(sum/float64(len(catUnits))*100) / 100,
			Min:     math.Round(min*100) / 100,
			Max:     math.Round(max*100) / 100,
		})
	}

	return result
}

// GetPerformanceDistribution returns the distribution of performance levels
func (r *Repository) GetPerformanceDistribution(units []entity.GRCUnit) []entity.PerformanceLevel {
	levels := []entity.PerformanceLevel{
		{Level: "Excellent", Min: 90, Max: 100, Count: 0},
		{Level: "Good", Min: 80, Max: 89, Count: 0},
		{Level: "Fair", Min: 70, Max: 79, Count: 0},
		{Level: "Needs Improvement", Min: 0, Max: 69, Count: 0},
	}

	for _, unit := range units {
		avg := unit.Average
		switch {
		case avg >= 90:
			levels[0].Count++
		case avg >= 80:
			levels[1].Count++
		case avg >= 70:
			levels[2].Count++
		default:
			levels[3].Count++
		}
	}

	return levels
}

// GetGapAnalysis returns gap analysis for a unit
func (r *Repository) GetGapAnalysis(unit *entity.GRCUnit, avg entity.GRCUnit) []entity.GapAnalysis {
	components := []struct {
		code    string
		unitVal float64
		avgVal  float64
	}{
		{"PSKI", unit.PSKI, avg.PSKI},
		{"PSPIP", unit.PSPIP, avg.PSPIP},
		{"PMRI", unit.PMRI, avg.PMRI},
		{"PIEPK", unit.PIEPK, avg.PIEPK},
		{"PWBKRB", unit.PWBKRB, avg.PWBKRB},
		{"PSAKIP", unit.PSAKIP, avg.PSAKIP},
	}

	var result []entity.GapAnalysis
	for _, comp := range components {
		gap := math.Round((comp.unitVal-comp.avgVal)*100) / 100
		status := "equal"
		if gap > 0 {
			status = "above"
		} else if gap < 0 {
			status = "below"
		}

		result = append(result, entity.GapAnalysis{
			Component: comp.code,
			Value:     comp.unitVal,
			Average:   comp.avgVal,
			Gap:       gap,
			Status:    status,
		})
	}

	return result
}

// GetWeakness returns the weakest component for a unit
func (r *Repository) GetWeakness(unit *entity.GRCUnit, avg entity.GRCUnit) entity.WeaknessPriority {
	components := []struct {
		code    string
		unitVal float64
		avgVal  float64
	}{
		{"PSKI", unit.PSKI, avg.PSKI},
		{"PSPIP", unit.PSPIP, avg.PSPIP},
		{"PMRI", unit.PMRI, avg.PMRI},
		{"PIEPK", unit.PIEPK, avg.PIEPK},
		{"PWBKRB", unit.PWBKRB, avg.PWBKRB},
		{"PSAKIP", unit.PSAKIP, avg.PSAKIP},
	}

	minIdx := 0
	for i, comp := range components {
		if comp.unitVal < components[minIdx].unitVal {
			minIdx = i
		}
	}

	weakness := components[minIdx]
	gap := weakness.unitVal - weakness.avgVal
	priority := "low"
	if gap < -10 {
		priority = "high"
	} else if gap < -5 {
		priority = "medium"
	}

	return entity.WeaknessPriority{
		Component:      weakness.code,
		Value:          weakness.unitVal,
		GapFromAverage: math.Round(gap*100) / 100,
		Priority:       priority,
	}
}

// GetStrength returns the strongest component for a unit
func (r *Repository) GetStrength(unit *entity.GRCUnit, avg entity.GRCUnit) entity.WeaknessPriority {
	components := []struct {
		code    string
		unitVal float64
		avgVal  float64
	}{
		{"PSKI", unit.PSKI, avg.PSKI},
		{"PSPIP", unit.PSPIP, avg.PSPIP},
		{"PMRI", unit.PMRI, avg.PMRI},
		{"PIEPK", unit.PIEPK, avg.PIEPK},
		{"PWBKRB", unit.PWBKRB, avg.PWBKRB},
		{"PSAKIP", unit.PSAKIP, avg.PSAKIP},
	}

	maxIdx := 0
	for i, comp := range components {
		if comp.unitVal > components[maxIdx].unitVal {
			maxIdx = i
		}
	}

	strength := components[maxIdx]
	gap := strength.unitVal - strength.avgVal

	return entity.WeaknessPriority{
		Component:      strength.code,
		Value:          strength.unitVal,
		GapFromAverage: math.Round(gap*100) / 100,
		Priority:       "strength",
	}
}

// GetRadarData returns radar chart data for a unit
func (r *Repository) GetRadarData(unit *entity.GRCUnit, avg entity.GRCUnit) entity.RadarData {
	return entity.RadarData{
		Labels:        []string{"PSKI", "PSPIP", "PMRI", "PIEPK", "PWBKRB", "PSAKIP"},
		Values:        []float64{unit.PSKI, unit.PSPIP, unit.PMRI, unit.PIEPK, unit.PWBKRB, unit.PSAKIP},
		AverageValues: []float64{avg.PSKI, avg.PSPIP, avg.PMRI, avg.PIEPK, avg.PWBKRB, avg.PSAKIP},
	}
}

// GetUnitsByIDs returns units by their IDs
func (r *Repository) GetUnitsByIDs(units []entity.GRCUnit, ids []int) []entity.GRCUnit {
	var result []entity.GRCUnit
	for _, id := range ids {
		unit := r.GetUnitByID(units, id)
		if unit != nil {
			result = append(result, *unit)
		}
	}
	return result
}
