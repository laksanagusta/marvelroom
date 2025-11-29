package handler

import (
	"context"

	vaccineUC "sandbox/internal/usecase/vaccine"
	"sandbox/pkg/pagination"

	"github.com/gofiber/fiber/v2"
)

type VaccineHandler struct {
	listMasterVaccinesUseCase    *vaccineUC.ListMasterVaccinesUseCase
	listCountriesUseCase         *vaccineUC.ListCountriesUseCase
	getCDCRecommendationsUseCase *vaccineUC.GetCDCRecommendationsUseCase
}

func NewVaccineHandler(
	listMasterVaccinesUseCase *vaccineUC.ListMasterVaccinesUseCase,
	listCountriesUseCase *vaccineUC.ListCountriesUseCase,
	getCDCRecommendationsUseCase *vaccineUC.GetCDCRecommendationsUseCase,
) *VaccineHandler {
	return &VaccineHandler{
		listMasterVaccinesUseCase:    listMasterVaccinesUseCase,
		listCountriesUseCase:         listCountriesUseCase,
		getCDCRecommendationsUseCase: getCDCRecommendationsUseCase,
	}
}

// ListMasterVaccines handles GET /api/v1/vaccines/vaccines
func (h *VaccineHandler) ListMasterVaccines(c *fiber.Ctx) error {
	// Parse query parameters like business trip handler
	queryParams := make(map[string]string)
	c.Context().QueryArgs().VisitAll(func(key, value []byte) {
		queryParams[string(key)] = string(value)
	})

	queryParser := &pagination.QueryParser{}
	params, err := queryParser.Parse(queryParams)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters: " + err.Error(),
		})
	}

	// Create request with parsed params
	req := vaccineUC.ListMasterVaccinesRequest{
		QueryParams: params,
	}

	// Get context from fiber
	ctx := c.Context()

	response, err := h.listMasterVaccinesUseCase.Execute(ctx, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get master vaccines",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Master vaccines retrieved successfully",
		"data":    response,
	})
}

// ListCountries handles GET /api/v1/vaccines/countries
func (h *VaccineHandler) ListCountries(c *fiber.Ctx) error {
	// Parse query parameters like business trip handler
	queryParams := make(map[string]string)
	c.Context().QueryArgs().VisitAll(func(key, value []byte) {
		queryParams[string(key)] = string(value)
	})

	queryParser := &pagination.QueryParser{}
	params, err := queryParser.Parse(queryParams)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters: " + err.Error(),
		})
	}

	req := vaccineUC.ListCountriesRequest{
		QueryParams: params,
	}
	response, err := h.listCountriesUseCase.Execute(context.Background(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":        response.Data,
		"message":     response.Message,
		"success":     response.Success,
		"page":        response.Page,
		"limit":       response.Limit,
		"total_items": response.TotalItems,
		"total_pages": response.TotalPages,
	})
}

// GetVaccineRecommendations handles GET /api/v1/vaccines/recommendations/:countryCode
func (h *VaccineHandler) GetVaccineRecommendations(c *fiber.Ctx) error {
	var req vaccineUC.GetCDCRecommendationsRequest

	// Parse country code from URL parameter
	req.CountryCode = c.Params("countryCode")
	if req.CountryCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Country code is required",
		})
	}

	// Parse query parameters
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid query parameters",
			"error":   err.Error(),
		})
	}

	// Get context from fiber
	ctx := c.Context()

	response, err := h.getCDCRecommendationsUseCase.Execute(ctx, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get vaccine recommendations",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Vaccine recommendations retrieved successfully",
		"data":    response,
	})
}
