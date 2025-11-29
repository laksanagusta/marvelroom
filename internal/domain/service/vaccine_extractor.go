package service

import "context"

// VaccineExtractor interface for extracting vaccine information from HTML
type VaccineExtractor interface {
	ExtractVaccineRecommendations(ctx context.Context, htmlContent string) (map[string]interface{}, error)
}
