package grc

import (
	"math"
	"sort"
	"time"

	"sandbox/internal/domain/entity"
	grcRepo "sandbox/internal/infrastructure/grc"
)

// GetOverviewUseCase handles getting dashboard overview
type GetOverviewUseCase struct {
	repo *grcRepo.Repository
}

// NewGetOverviewUseCase creates a new use case
func NewGetOverviewUseCase(repo *grcRepo.Repository) *GetOverviewUseCase {
	return &GetOverviewUseCase{repo: repo}
}

// Execute returns the dashboard overview
func (uc *GetOverviewUseCase) Execute() (*OverviewResponse, error) {
	// Fetch data once
	units, _, err := uc.repo.FetchData()
	if err != nil {
		return nil, err
	}

	// Calculate statistics
	avg, median, stdDev := uc.repo.CalculateStatistics(units)

	// Get top performers
	topPerformers := uc.repo.GetTopPerformers(units, 5)
	var topSummary []UnitSummary
	for _, u := range topPerformers {
		rank := uc.repo.GetUnitRank(units, u.ID)
		topSummary = append(topSummary, UnitSummary{
			ID:       u.ID,
			Name:     u.Name,
			Category: u.Category,
			Average:  u.Average,
			Rank:     rank,
		})
	}

	// Get bottom performers
	bottomPerformers := uc.repo.GetBottomPerformers(units, 5)
	var bottomSummary []UnitSummary
	for _, u := range bottomPerformers {
		rank := uc.repo.GetUnitRank(units, u.ID)
		bottomSummary = append(bottomSummary, UnitSummary{
			ID:       u.ID,
			Name:     u.Name,
			Category: u.Category,
			Average:  u.Average,
			Rank:     rank,
		})
	}

	// Get component stats
	componentStats := uc.repo.GetComponentStats(units)
	sort.Slice(componentStats, func(i, j int) bool {
		return componentStats[i].UnitsBelow80 > componentStats[j].UnitsBelow80
	})

	// Get category stats
	categoryStats := uc.repo.GetCategoryStats(units)

	// Get performance distribution
	perfDist := uc.repo.GetPerformanceDistribution(units)

	return &OverviewResponse{
		Statistics: OverviewStatistics{
			TotalUnits:   len(units),
			AverageScore: avg,
			Median:       median,
			StdDeviation: stdDev,
		},
		TopPerformers:           topSummary,
		BottomPerformers:        bottomSummary,
		WeakestComponents:       componentStats,
		CategorySummary:         categoryStats,
		PerformanceDistribution: perfDist,
		LastUpdated:             time.Now().Format(time.RFC3339),
	}, nil
}

// ListUnitsUseCase handles listing all units
type ListUnitsUseCase struct {
	repo *grcRepo.Repository
}

// NewListUnitsUseCase creates a new use case
func NewListUnitsUseCase(repo *grcRepo.Repository) *ListUnitsUseCase {
	return &ListUnitsUseCase{repo: repo}
}

// Execute returns all units
func (uc *ListUnitsUseCase) Execute(category string, sortBy string, ascending bool) (*UnitListResponse, error) {
	// Fetch data once
	allUnits, _, err := uc.repo.FetchData()
	if err != nil {
		return nil, err
	}

	// Use allUnits for filtering, keep original for rank calculation
	units := allUnits

	// Filter by category if specified
	if category != "" {
		units = uc.repo.GetUnitsByCategory(allUnits, category)
	}

	// Sort units
	switch sortBy {
	case "name":
		sort.Slice(units, func(i, j int) bool {
			if ascending {
				return units[i].Name < units[j].Name
			}
			return units[i].Name > units[j].Name
		})
	case "average":
		sort.Slice(units, func(i, j int) bool {
			if ascending {
				return units[i].Average < units[j].Average
			}
			return units[i].Average > units[j].Average
		})
	default:
		sort.Slice(units, func(i, j int) bool {
			if ascending {
				return units[i].ID < units[j].ID
			}
			return units[i].ID > units[j].ID
		})
	}

	var items []UnitListItem
	for _, u := range units {
		rank := uc.repo.GetUnitRank(allUnits, u.ID)
		percentile := uc.repo.GetPercentile(allUnits, u.ID)
		items = append(items, UnitListItem{
			ID:         u.ID,
			Name:       u.Name,
			Category:   u.Category,
			Average:    u.Average,
			Rank:       rank,
			Percentile: percentile,
			Scores: UnitScores{
				PSKI:   u.PSKI,
				PSPIP:  u.PSPIP,
				PMRI:   u.PMRI,
				PIEPK:  u.PIEPK,
				PWBKRB: u.PWBKRB,
				PSAKIP: u.PSAKIP,
			},
		})
	}

	return &UnitListResponse{
		Units:      items,
		TotalCount: len(items),
	}, nil
}

// GRCUnit alias for entity
type GRCUnit = entity.GRCUnit

// GetUnitDetailUseCase handles getting unit detail
type GetUnitDetailUseCase struct {
	repo *grcRepo.Repository
}

// NewGetUnitDetailUseCase creates a new use case
func NewGetUnitDetailUseCase(repo *grcRepo.Repository) *GetUnitDetailUseCase {
	return &GetUnitDetailUseCase{repo: repo}
}

// Execute returns unit detail
func (uc *GetUnitDetailUseCase) Execute(id int) (*UnitDetailResponse, error) {
	// Fetch data once
	units, avg, err := uc.repo.FetchData()
	if err != nil {
		return nil, err
	}

	unit := uc.repo.GetUnitByID(units, id)
	if unit == nil {
		return nil, nil
	}

	rank := uc.repo.GetUnitRank(units, unit.ID)
	percentile := uc.repo.GetPercentile(units, unit.ID)

	radarData := uc.repo.GetRadarData(unit, avg)
	gapAnalysis := uc.repo.GetGapAnalysis(unit, avg)
	weakness := uc.repo.GetWeakness(unit, avg)
	strength := uc.repo.GetStrength(unit, avg)

	// Calculate category comparison
	categoryUnits := uc.repo.GetUnitsByCategory(units, unit.Category)
	var categorySum float64
	for _, u := range categoryUnits {
		categorySum += u.Average
	}
	categoryAvg := categorySum / float64(len(categoryUnits))
	categoryAvg = math.Round(categoryAvg*100) / 100

	return &UnitDetailResponse{
		Unit: UnitInfo{
			ID:         unit.ID,
			Name:       unit.Name,
			Category:   unit.Category,
			Average:    unit.Average,
			Rank:       rank,
			Percentile: percentile,
		},
		Scores: UnitScores{
			PSKI:   unit.PSKI,
			PSPIP:  unit.PSPIP,
			PMRI:   unit.PMRI,
			PIEPK:  unit.PIEPK,
			PWBKRB: unit.PWBKRB,
			PSAKIP: unit.PSAKIP,
		},
		RadarData:   radarData,
		GapAnalysis: gapAnalysis,
		Weakness:    weakness,
		Strength:    strength,
		CategoryComparison: &CategoryComparison{
			CategoryName:    unit.Category,
			CategoryAverage: categoryAvg,
			GapToCategory:   math.Round((unit.Average-categoryAvg)*100) / 100,
		},
	}, nil
}

// CompareUnitsUseCase handles comparing multiple units
type CompareUnitsUseCase struct {
	repo *grcRepo.Repository
}

// NewCompareUnitsUseCase creates a new use case
func NewCompareUnitsUseCase(repo *grcRepo.Repository) *CompareUnitsUseCase {
	return &CompareUnitsUseCase{repo: repo}
}

// Execute returns comparison data
func (uc *CompareUnitsUseCase) Execute(ids []int) (*CompareResponse, error) {
	// Fetch data once
	allUnits, avg, err := uc.repo.FetchData()
	if err != nil {
		return nil, err
	}

	units := uc.repo.GetUnitsByIDs(allUnits, ids)

	var compareUnits []CompareUnit
	for _, u := range units {
		rank := uc.repo.GetUnitRank(allUnits, u.ID)
		compareUnits = append(compareUnits, CompareUnit{
			ID:       u.ID,
			Name:     u.Name,
			Category: u.Category,
			Average:  u.Average,
			Rank:     rank,
			Values:   []float64{u.PSKI, u.PSPIP, u.PMRI, u.PIEPK, u.PWBKRB, u.PSAKIP},
		})
	}

	return &CompareResponse{
		Units:       compareUnits,
		RadarLabels: []string{"PSKI", "PSPIP", "PMRI", "PIEPK", "PWBKRB", "PSAKIP"},
		Average: entity.RadarData{
			Labels:        []string{"PSKI", "PSPIP", "PMRI", "PIEPK", "PWBKRB", "PSAKIP"},
			Values:        []float64{avg.PSKI, avg.PSPIP, avg.PMRI, avg.PIEPK, avg.PWBKRB, avg.PSAKIP},
			AverageValues: []float64{avg.PSKI, avg.PSPIP, avg.PMRI, avg.PIEPK, avg.PWBKRB, avg.PSAKIP},
		},
	}, nil
}

// GetCategoriesUseCase handles getting category breakdown
type GetCategoriesUseCase struct {
	repo *grcRepo.Repository
}

// NewGetCategoriesUseCase creates a new use case
func NewGetCategoriesUseCase(repo *grcRepo.Repository) *GetCategoriesUseCase {
	return &GetCategoriesUseCase{repo: repo}
}

// Execute returns category breakdown
func (uc *GetCategoriesUseCase) Execute() (*CategoryResponse, error) {
	// Fetch data once
	units, _, err := uc.repo.FetchData()
	if err != nil {
		return nil, err
	}

	stats := uc.repo.GetCategoryStats(units)

	var categories []CategoryDetail

	for _, stat := range stats {
		catUnits := uc.repo.GetUnitsByCategory(units, stat.Name)
		if len(catUnits) == 0 {
			continue
		}

		sort.Slice(catUnits, func(i, j int) bool {
			return catUnits[i].Average > catUnits[j].Average
		})

		topUnit := catUnits[0]
		bottomUnit := catUnits[len(catUnits)-1]
		topRank := uc.repo.GetUnitRank(units, topUnit.ID)
		bottomRank := uc.repo.GetUnitRank(units, bottomUnit.ID)

		categories = append(categories, CategoryDetail{
			Name:    stat.Name,
			Count:   stat.Count,
			Average: stat.Average,
			Min:     stat.Min,
			Max:     stat.Max,
			TopUnit: &UnitSummary{
				ID:       topUnit.ID,
				Name:     topUnit.Name,
				Category: topUnit.Category,
				Average:  topUnit.Average,
				Rank:     topRank,
			},
			BottomUnit: &UnitSummary{
				ID:       bottomUnit.ID,
				Name:     bottomUnit.Name,
				Category: bottomUnit.Category,
				Average:  bottomUnit.Average,
				Rank:     bottomRank,
			},
		})
	}

	return &CategoryResponse{Categories: categories}, nil
}
