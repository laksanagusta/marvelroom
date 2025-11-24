package handler

import (
	"context"

	"sandbox/internal/usecase/business_trip"

	"github.com/gofiber/fiber/v2"
)

type AssigneeHandler struct {
	addAssigneeUseCase     *business_trip.AddAssigneeUseCase
	getAssigneeUseCase     *business_trip.GetAssigneeUseCase
	updateAssigneeUseCase  *business_trip.UpdateAssigneeUseCase
	deleteAssigneeUseCase  *business_trip.DeleteAssigneeUseCase
	listAssigneesUseCase   *business_trip.ListAssigneesUseCase
}

func NewAssigneeHandler(
	addAssigneeUseCase *business_trip.AddAssigneeUseCase,
	getAssigneeUseCase *business_trip.GetAssigneeUseCase,
	updateAssigneeUseCase *business_trip.UpdateAssigneeUseCase,
	deleteAssigneeUseCase *business_trip.DeleteAssigneeUseCase,
	listAssigneesUseCase *business_trip.ListAssigneesUseCase,
) *AssigneeHandler {
	return &AssigneeHandler{
		addAssigneeUseCase:    addAssigneeUseCase,
		getAssigneeUseCase:    getAssigneeUseCase,
		updateAssigneeUseCase: updateAssigneeUseCase,
		deleteAssigneeUseCase: deleteAssigneeUseCase,
		listAssigneesUseCase:  listAssigneesUseCase,
	}
}

// CreateAssignee creates a new assignee for a business trip
func (h *AssigneeHandler) CreateAssignee(c *fiber.Ctx) error {
	tripID := c.Params("tripId")
	if tripID == "" {
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

	_, err := h.addAssigneeUseCase.Execute(context.Background(), tripID, &req)
	if err != nil {
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

// ListAssignees lists all assignees for a business trip
func (h *AssigneeHandler) ListAssignees(c *fiber.Ctx) error {
	tripID := c.Params("tripId")
	if tripID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Business trip ID is required",
		})
	}

	response, err := h.listAssigneesUseCase.Execute(context.Background(), tripID)
	if err != nil {
		if err != nil && err.Error() == "business trip not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Business trip not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get assignees",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Assignees retrieved successfully",
		"data":    response,
	})
}

// GetAssignee gets a specific assignee by ID
func (h *AssigneeHandler) GetAssignee(c *fiber.Ctx) error {
	assigneeID := c.Params("assigneeId")
	if assigneeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Assignee ID is required",
		})
	}

	response, err := h.getAssigneeUseCase.Execute(context.Background(), assigneeID)
	if err != nil {
		if err != nil && err.Error() == "assignee not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Assignee not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get assignee",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Assignee retrieved successfully",
		"data":    response,
	})
}

// UpdateAssignee updates a specific assignee
func (h *AssigneeHandler) UpdateAssignee(c *fiber.Ctx) error {
	tripId := c.Params("tripId")
	assigneeID := c.Params("assigneeId")
	if tripId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Business trip ID is required",
		})
	}
	if assigneeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Assignee ID is required",
		})
	}

	var req business_trip.UpdateAssigneeRequest

	// Parse path parameters first
	req.BusinessTripID = tripId
	req.AssigneeID = assigneeID

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

	_, err := h.updateAssigneeUseCase.Execute(context.Background(), req)
	if err != nil {
		if err != nil && err.Error() == "assignee not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Assignee not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update assignee",
			"details": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// DeleteAssignee deletes a specific assignee with parent validation
func (h *AssigneeHandler) DeleteAssignee(c *fiber.Ctx) error {
	tripId := c.Params("tripId")
	assigneeID := c.Params("assigneeId")
	if tripId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Business trip ID is required",
		})
	}
	if assigneeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Assignee ID is required",
		})
	}

	// Manual parent validation before deleting
	// Get assignee to verify it belongs to business trip
	assignee, err := h.getAssigneeUseCase.Execute(context.Background(), assigneeID)
	if err != nil {
		if err != nil && err.Error() == "assignee not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Assignee not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get assignee",
			"details": err.Error(),
		})
	}

	if assignee.BusinessTripID != tripId {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Assignee does not belong to the specified business trip",
		})
	}

	err = h.deleteAssigneeUseCase.Execute(context.Background(), assigneeID)
	if err != nil {
		if err != nil && err.Error() == "assignee not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Assignee not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete assignee",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Assignee deleted successfully",
	})
}