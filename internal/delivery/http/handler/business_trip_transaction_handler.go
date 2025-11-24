package handler

import (
	"context"

	"sandbox/internal/usecase/business_trip"

	"github.com/gofiber/fiber/v2"
)

// BusinessTripTransactionHandler handles HTTP requests for business trip transactions
type BusinessTripTransactionHandler struct {
	addTransactionUseCase    *business_trip.AddTransactionUseCase
	getTransactionUseCase    *business_trip.GetTransactionUseCase
	updateTransactionUseCase *business_trip.UpdateTransactionUseCase
	deleteTransactionUseCase *business_trip.DeleteTransactionUseCase
	listTransactionsUseCase  *business_trip.ListTransactionsUseCase
	getAssigneeUseCase       *business_trip.GetAssigneeUseCase
}

func NewBusinessTripTransactionHandler(
	addTransactionUseCase *business_trip.AddTransactionUseCase,
	getTransactionUseCase *business_trip.GetTransactionUseCase,
	updateTransactionUseCase *business_trip.UpdateTransactionUseCase,
	deleteTransactionUseCase *business_trip.DeleteTransactionUseCase,
	listTransactionsUseCase *business_trip.ListTransactionsUseCase,
	getAssigneeUseCase *business_trip.GetAssigneeUseCase,
) *BusinessTripTransactionHandler {
	return &BusinessTripTransactionHandler{
		addTransactionUseCase:    addTransactionUseCase,
		getTransactionUseCase:    getTransactionUseCase,
		updateTransactionUseCase: updateTransactionUseCase,
		deleteTransactionUseCase: deleteTransactionUseCase,
		listTransactionsUseCase:  listTransactionsUseCase,
		getAssigneeUseCase:       getAssigneeUseCase,
	}
}

// Create creates a new transaction for an assignee
func (h *BusinessTripTransactionHandler) Create(c *fiber.Ctx) error {
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

// List lists all transactions for an assignee
func (h *BusinessTripTransactionHandler) List(c *fiber.Ctx) error {
	assigneeID := c.Params("assigneeId")
	if assigneeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Assignee ID is required",
		})
	}

	response, err := h.listTransactionsUseCase.Execute(context.Background(), assigneeID)
	if err != nil {
		if err != nil && err.Error() == "assignee not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Assignee not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get transactions",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Transactions retrieved successfully",
		"data":    response,
	})
}

// Get gets a specific transaction by ID
func (h *BusinessTripTransactionHandler) Get(c *fiber.Ctx) error {
	transactionID := c.Params("transactionId")
	if transactionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Transaction ID is required",
		})
	}

	response, err := h.getTransactionUseCase.Execute(context.Background(), transactionID)
	if err != nil {
		if err != nil && err.Error() == "transaction not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Transaction not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get transaction",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Transaction retrieved successfully",
		"data":    response,
	})
}

// Update updates a specific transaction
func (h *BusinessTripTransactionHandler) Update(c *fiber.Ctx) error {
	tripId := c.Params("tripId")
	assigneeID := c.Params("assigneeId")
	transactionID := c.Params("transactionId")
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
	if transactionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Transaction ID is required",
		})
	}

	var req business_trip.UpdateTransactionRequest

	// Parse path parameters first
	req.BusinessTripID = tripId
	req.AssigneeID = assigneeID
	req.TransactionID = transactionID

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

	_, err := h.updateTransactionUseCase.Execute(context.Background(), req)
	if err != nil {
		if err != nil && err.Error() == "transaction not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Transaction not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update transaction",
			"details": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// Delete deletes a specific transaction with parent validation
func (h *BusinessTripTransactionHandler) Delete(c *fiber.Ctx) error {
	tripId := c.Params("tripId")
	assigneeID := c.Params("assigneeId")
	transactionID := c.Params("transactionId")
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
	if transactionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Transaction ID is required",
		})
	}

	// Manual parent validation before deleting
	// Verify transaction belongs to the assignee and assignee belongs to business trip
	transaction, err := h.getTransactionUseCase.Execute(context.Background(), transactionID)
	if err != nil {
		if err != nil && err.Error() == "transaction not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Transaction not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get transaction",
			"details": err.Error(),
		})
	}

	if transaction.AssigneeID != assigneeID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Transaction does not belong to the specified assignee",
		})
	}

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

	// Delete transaction
	err = h.deleteTransactionUseCase.Execute(context.Background(), transactionID)
	if err != nil {
		if err != nil && err.Error() == "transaction not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Transaction not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete transaction",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Transaction deleted successfully",
	})
}
