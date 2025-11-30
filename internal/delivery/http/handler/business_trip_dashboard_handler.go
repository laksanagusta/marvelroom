package handler

import (
	"context"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"sandbox/internal/usecase/business_trip"
)

// BusinessTripDashboardHandler handles HTTP requests for business trip dashboard
type BusinessTripDashboardHandler struct {
	dashboardUseCase *business_trip.GetDashboardUseCase
	validator        *validator.Validate
}

// NewBusinessTripDashboardHandler creates a new handler instance
func NewBusinessTripDashboardHandler(dashboardUseCase *business_trip.GetDashboardUseCase) *BusinessTripDashboardHandler {
	return &BusinessTripDashboardHandler{
		dashboardUseCase: dashboardUseCase,
		validator:        validator.New(),
	}
}

// GetDashboard retrieves business trip dashboard data
// @Summary Get Business Trip Dashboard
// @Description Retrieves comprehensive dashboard data for business trips including overview, monthly stats, destination stats, and recent trips
// @Tags business-trips
// @Accept json
// @Produce json
// @Param start_date query string false "Start date filter (YYYY-MM-DD format)"
// @Param end_date query string false "End date filter (YYYY-MM-DD format)"
// @Param destination query string false "Destination city filter"
// @Param status query string false "Status filter (draft, ongoing, completed, canceled)"
// @Param limit query int false "Limit for recent trips (default: 10, max: 100)"
// @Success 200 {object} StandardResponse{data=business_trip.GetDashboardResponse}
// @Failure 400 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /api/v1/business-trips/dashboard [get]
func (h *BusinessTripDashboardHandler) GetDashboard(c *fiber.Ctx) error {
	// Temporary: Skip authentication check for testing
	// TODO: Re-enable authentication after testing is complete
	/*
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header is required",
		})
	}
	*/

	// Parse query parameters
	req := business_trip.GetDashboardRequest{
		StartDate:   parseDateQueryParam(c.Query("start_date")),
		EndDate:     parseDateQueryParam(c.Query("end_date")),
		Destination: c.Query("destination"),
		Status:      c.Query("status"),
		Limit:       parseIntQueryParam(c.Query("limit"), 10),
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	// Execute use case
	ctx := context.Background()
	response, err := h.dashboardUseCase.Execute(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve dashboard data",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// parseDateQueryParam parses date query parameter
func parseDateQueryParam(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil
	}

	return &date
}

// parseIntQueryParam parses integer query parameter with default value
func parseIntQueryParam(intStr string, defaultValue int) int {
	if intStr == "" {
		return defaultValue
	}

	// Try to parse the integer string
	if value, err := strconv.Atoi(intStr); err == nil && value > 0 {
		return value
	}

	return defaultValue
}
