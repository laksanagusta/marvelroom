package handler

import (
	"sandbox/internal/delivery/http/middleware"
	"sandbox/internal/usecase/business_trip"
	"sandbox/pkg/pagination"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// BusinessTripVerificationHandler handles HTTP requests for business trip verification
type BusinessTripVerificationHandler struct {
	verifyUseCase           *business_trip.VerifyBusinessTripUseCase
	listVerificatorsUseCase *business_trip.ListVerificatorsUseCase
	validator               *validator.Validate
}

// NewBusinessTripVerificationHandler creates a new handler instance
func NewBusinessTripVerificationHandler(
	verifyUseCase *business_trip.VerifyBusinessTripUseCase,
	listVerificatorsUseCase *business_trip.ListVerificatorsUseCase,
) *BusinessTripVerificationHandler {
	return &BusinessTripVerificationHandler{
		verifyUseCase:           verifyUseCase,
		listVerificatorsUseCase: listVerificatorsUseCase,
		validator:               validator.New(),
	}
}

// VerifyBusinessTrip handles the verification of a business trip
// @Summary Verify Business Trip
// @Description Allows a verificator to approve or reject a business trip that is in ready_to_verify status
// @Tags business-trips
// @Accept json
// @Produce json
// @Param tripId path string true "Business Trip ID"
// @Param request body business_trip.VerifyBusinessTripRequest true "Verification Request"
// @Success 200 {object} StandardResponse{data=business_trip.VerifyBusinessTripResponse}
// @Failure 400 {object} StandardResponse
// @Failure 401 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /api/v1/business-trips/{tripId}/verify [post]
func (h *BusinessTripVerificationHandler) VerifyBusinessTrip(c *fiber.Ctx) error {
	// Get business trip ID from URL parameters
	businessTripID := c.Params("tripId")
	if businessTripID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Business trip ID is required",
		})
	}

	// Parse request body
	var req business_trip.VerifyBusinessTripRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Set business trip ID from URL parameter
	req.BusinessTripID = businessTripID

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	// Additional validation using the request's Validate method
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	authenticatedUser, _ := middleware.GetAuthenticatedUser(c)

	// Execute use case
	// The context should contain user_id from authentication middleware
	response, err := h.verifyUseCase.Execute(c.Context(), req, *authenticatedUser)
	if err != nil {
		// Handle authentication error
		if err.Error() == "authentication error: user not authenticated or user_id not found in context" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Authentication required",
			})
		}

		// Handle different types of errors appropriately
		if err.Error() == "business trip must be in ready_to_verify status to be verified" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		// Check if it's a "not found" error
		if err.Error() == "failed to get business trip: record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Business trip not found",
			})
		}

		// Check if it's a verificator not found or unauthorized error
		if err.Error() == "failed to get verificator: record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "You are not assigned as a verificator for this business trip",
			})
		}

		if err.Error() == "verificator has already approved this business trip" ||
			err.Error() == "verificator has already rejected this business trip" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		// Generic server error for other cases
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to verify business trip",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// ListVerificators lists business trip verificators with pagination and filtering
// @Summary List Business Trip Verificators
// @Description Retrieves a paginated list of business trip verificators with filtering and sorting capabilities
// @Tags business-trips
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20, max: 100)"
// @Param sort query string false "Sort fields (e.g., 'status asc,business_trip_number desc')"
// @Param status query string false "Filter by verification status (pending, approved, rejected)"
// @Param business_trip_status query string false "Filter by business trip status (draft, ongoing, completed, canceled, ready_to_verify)"
// @Param user_id query string false "Filter by user ID"
// @Param destination_city query string false "Filter by destination city"
// @Param activity_purpose query string false "Filter by activity purpose (contains)"
// @Success 200 {object} pagination.PagedResponse{data=[]business_trip.ListVerificatorsResponse}
// @Failure 400 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /api/v1/business-trips/verificators [get]
func (h *BusinessTripVerificationHandler) ListVerificators(c *fiber.Ctx) error {
	// Parse query parameters using the same pattern as ListBusinessTrips
	queryParams := make(map[string]string)
	c.Context().QueryArgs().VisitAll(func(key, value []byte) {
		queryParams[string(key)] = string(value)
	})

	queryParser := &pagination.QueryParser{}
	params, err := queryParser.Parse(queryParams)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid query parameters: " + err.Error(),
		})
	}

	// Execute use case
	verificators, pagination, err := h.listVerificatorsUseCase.Execute(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to retrieve verificators",
			"details": err.Error(),
		})
	}

	// Set the data in pagination response
	pagination.Data = verificators

	return c.JSON(pagination)
}
