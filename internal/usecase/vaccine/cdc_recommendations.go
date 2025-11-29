package vaccine

import (
	"context"
	"fmt"
	"time"

	"sandbox/internal/domain/service"
)

type GetCDCRecommendationsRequest struct {
	CountryCode string `param:"countryCode"`
	Language    string `query:"language"`
}

type GetCDCRecommendationsResponse struct {
	Data        *service.CountryVaccineRecommendation `json:"data"`
	Message     string                                `json:"message"`
	Success     bool                                  `json:"success"`
	CountryCode string                                `json:"country_code"`
	Language    string                                `json:"language"`
	LastUpdated time.Time                             `json:"last_updated"`
}

type GetCDCRecommendationsUseCase struct {
	cdcService *service.CDCService
}

func NewGetCDCRecommendationsUseCase(cdcService *service.CDCService) *GetCDCRecommendationsUseCase {
	return &GetCDCRecommendationsUseCase{
		cdcService: cdcService,
	}
}

func (r GetCDCRecommendationsRequest) Validate() error {
	if r.CountryCode == "" {
		return fmt.Errorf("country code is required")
	}
	return nil
}

func (uc *GetCDCRecommendationsUseCase) Execute(ctx context.Context, req *GetCDCRecommendationsRequest) (*GetCDCRecommendationsResponse, error) {
	// Default language to English if not provided
	language := req.Language
	if language == "" {
		language = "en"
	}

	recommendation, err := uc.cdcService.GetVaccineRecommendationsByCountry(ctx, req.CountryCode, language)
	if err != nil {
		return nil, err
	}

	return &GetCDCRecommendationsResponse{
		Data:        recommendation,
		Message:     "CDC vaccine recommendations retrieved successfully",
		Success:     true,
		CountryCode: req.CountryCode,
		Language:    language,
		LastUpdated: time.Now(),
	}, nil
}
