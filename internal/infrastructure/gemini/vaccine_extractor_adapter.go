package gemini

import (
	"context"

	"sandbox/internal/domain/service"
)

// VaccineExtractorAdapter wraps the Gemini client to implement VaccineExtractor interface
type VaccineExtractorAdapter struct {
	client *Client
}

// NewVaccineExtractorAdapter creates a new adapter for vaccine extraction
func NewVaccineExtractorAdapter(client *Client) service.VaccineExtractor {
	return &VaccineExtractorAdapter{
		client: client,
	}
}

// ExtractVaccineRecommendations extracts vaccine information from HTML using Gemini
func (a *VaccineExtractorAdapter) ExtractVaccineRecommendations(ctx context.Context, htmlContent string) (map[string]interface{}, error) {
	return a.client.ExtractVaccineRecommendations(ctx, htmlContent)
}