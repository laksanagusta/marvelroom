package handler

import (
	"context"

	"sandbox/internal/delivery/http/middleware"
	"sandbox/internal/usecase/business_trip"
	"sandbox/pkg/pagination"

	"github.com/gofiber/fiber/v2"
)

type BusinessTripHandler struct {
	createBusinessTripUseCase              *business_trip.CreateBusinessTripUseCase
	getBusinessTripUseCase                 *business_trip.GetBusinessTripUseCase
	updateBusinessTripUseCase              *business_trip.UpdateBusinessTripUseCase
	updateBusinessTripWithAssigneesUseCase *business_trip.UpdateBusinessTripWithAssigneesUseCase
	deleteBusinessTripUseCase              *business_trip.DeleteBusinessTripUseCase
	listBusinessTripsUseCase               *business_trip.ListBusinessTripsUseCase
	addAssigneeUseCase                     *business_trip.AddAssigneeUseCase
	addTransactionUseCase                  *business_trip.AddTransactionUseCase
	getBusinessTripSummaryUseCase          *business_trip.GetBusinessTripSummaryUseCase
	getAssigneeSummaryUseCase              *business_trip.GetAssigneeSummaryUseCase
}

func NewBusinessTripHandler(
	createBusinessTripUseCase *business_trip.CreateBusinessTripUseCase,
	getBusinessTripUseCase *business_trip.GetBusinessTripUseCase,
	updateBusinessTripUseCase *business_trip.UpdateBusinessTripUseCase,
	updateBusinessTripWithAssigneesUseCase *business_trip.UpdateBusinessTripWithAssigneesUseCase,
	deleteBusinessTripUseCase *business_trip.DeleteBusinessTripUseCase,
	listBusinessTripsUseCase *business_trip.ListBusinessTripsUseCase,
	addAssigneeUseCase *business_trip.AddAssigneeUseCase,
	addTransactionUseCase *business_trip.AddTransactionUseCase,
	getBusinessTripSummaryUseCase *business_trip.GetBusinessTripSummaryUseCase,
	getAssigneeSummaryUseCase *business_trip.GetAssigneeSummaryUseCase,
) *BusinessTripHandler {
	return &BusinessTripHandler{
		createBusinessTripUseCase:              createBusinessTripUseCase,
		getBusinessTripUseCase:                 getBusinessTripUseCase,
		updateBusinessTripUseCase:              updateBusinessTripUseCase,
		updateBusinessTripWithAssigneesUseCase: updateBusinessTripWithAssigneesUseCase,
		deleteBusinessTripUseCase:              deleteBusinessTripUseCase,
		listBusinessTripsUseCase:               listBusinessTripsUseCase,
		addAssigneeUseCase:                     addAssigneeUseCase,
		addTransactionUseCase:                  addTransactionUseCase,
		getBusinessTripSummaryUseCase:          getBusinessTripSummaryUseCase,
		getAssigneeSummaryUseCase:              getAssigneeSummaryUseCase,
	}
}

// CreateBusinessTrip creates a new business trip
func (h *BusinessTripHandler) CreateBusinessTrip(c *fiber.Ctx) error {
	// Get authenticated user from context
	user, err := middleware.GetAuthenticatedUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Authentication required",
			"details": err.Error(),
		})
	}

	// Example: Use the authenticated user data
	_ = user // User data can be used for authorization or logging
	_ = user.ID
	_ = user.Username
	_ = user.GetFullName()
	_ = user.GetPrimaryRole()
	_ = user.Organization.Name

	var req business_trip.BusinessTripRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	// Call usecase directly
	response, err := h.createBusinessTripUseCase.Execute(context.Background(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create business trip",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Business trip created successfully",
		"data":    response,
	})
}

// GetBusinessTrip gets a business trip by ID
func (h *BusinessTripHandler) GetBusinessTrip(c *fiber.Ctx) error {
	id := c.Params("tripId")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Business trip ID is required",
		})
	}

	response, err := h.getBusinessTripUseCase.Execute(context.Background(), id)
	if err != nil {
		// Check if it's a not found error
		if err != nil && err.Error() == "business trip not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Business trip not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get business trip",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Business trip retrieved successfully",
		"data":    response,
	})
}

// UpdateBusinessTrip updates a business trip
func (h *BusinessTripHandler) UpdateBusinessTrip(c *fiber.Ctx) error {
	tripId := c.Params("tripId")
	if tripId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Business trip ID is required",
		})
	}

	var req business_trip.UpdateBusinessTripRequest

	// Parse path parameters first
	req.BusinessTripID = tripId

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	// Call usecase directly
	_, err := h.updateBusinessTripUseCase.Execute(context.Background(), req)
	if err != nil {
		if err != nil && (err.Error() == "business trip not found" || err.Error() == "entity.ErrBusinessTripNotFound") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Business trip not found",
			})
		}
		if err != nil && err.Error() == "invalid date range" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid date range",
				"details": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update business trip",
			"details": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// UpdateBusinessTripWithAssignees updates a business trip and replaces all its assignees and transactions
func (h *BusinessTripHandler) UpdateBusinessTripWithAssignees(c *fiber.Ctx) error {
	tripId := c.Params("tripId")
	if tripId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Business trip ID is required",
		})
	}

	var req business_trip.UpdateBusinessTripWithAssigneesRequest

	// Parse path parameters first
	req.BusinessTripID = tripId

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	// Call usecase directly
	_, err := h.updateBusinessTripWithAssigneesUseCase.Execute(context.Background(), req)
	if err != nil {
		if err != nil && (err.Error() == "business trip not found" || err.Error() == "entity.ErrBusinessTripNotFound") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Business trip not found",
			})
		}
		if err != nil && err.Error() == "invalid date range" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid date range",
				"details": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update business trip with assignees",
			"details": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// DeleteBusinessTrip deletes a business trip
func (h *BusinessTripHandler) DeleteBusinessTrip(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Business trip ID is required",
		})
	}

	err := h.deleteBusinessTripUseCase.Execute(context.Background(), id)
	if err != nil {
		// Check if it's a not found error
		if err != nil && err.Error() == "business trip not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Business trip not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete business trip",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Business trip deleted successfully",
	})
}

// ListBusinessTrips lists business trips with pagination and filtering
func (h *BusinessTripHandler) ListBusinessTrips(c *fiber.Ctx) error {
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

	businessTrips, pagination, err := h.listBusinessTripsUseCase.Execute(context.Background(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	pagination.Data = businessTrips

	return c.JSON(pagination)
}

// AddAssignee adds an assignee to a business trip
func (h *BusinessTripHandler) AddAssignee(c *fiber.Ctx) error {
	businessTripID := c.Params("businessTripId")
	if businessTripID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Business trip ID is required",
		})
	}

	var req business_trip.AssigneeRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	// Validate nested transactions
	for _, tx := range req.Transactions {
		if err := tx.Validate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Transaction validation failed",
				"details": err.Error(),
			})
		}
	}

	_, err := h.addAssigneeUseCase.Execute(context.Background(), businessTripID, &req)
	if err != nil {
		// Check if it's a not found error
		if err != nil && err.Error() == "business trip not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Business trip not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to add assignee",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).SendStatus(fiber.StatusCreated)
}

// AddTransaction adds a transaction to an assignee
func (h *BusinessTripHandler) AddTransaction(c *fiber.Ctx) error {
	assigneeID := c.Params("assigneeId")
	if assigneeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Assignee ID is required",
		})
	}

	var req business_trip.TransactionRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	_, err := h.addTransactionUseCase.Execute(context.Background(), assigneeID, req)
	if err != nil {
		if err != nil && err.Error() == "assignee not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Assignee not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to add transaction",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).SendStatus(fiber.StatusCreated)
}

// GetBusinessTripSummary gets a summary of a business trip
func (h *BusinessTripHandler) GetBusinessTripSummary(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Business trip ID is required",
		})
	}

	summary, err := h.getBusinessTripSummaryUseCase.Execute(context.Background(), id)
	if err != nil {
		// Check if it's a not found error
		if err != nil && err.Error() == "business trip not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Business trip not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get business trip summary",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Business trip summary retrieved successfully",
		"data":    summary,
	})
}

// GetAssigneeSummary gets a summary of an assignee
func (h *BusinessTripHandler) GetAssigneeSummary(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Assignee ID is required",
		})
	}

	summary, err := h.getAssigneeSummaryUseCase.Execute(context.Background(), id)
	if err != nil {
		if err != nil && err.Error() == "assignee not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Assignee not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get assignee summary",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Assignee summary retrieved successfully",
		"data":    summary,
	})
}
