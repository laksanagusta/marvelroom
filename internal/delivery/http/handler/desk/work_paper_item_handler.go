package desk

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"sandbox/internal/usecase/work_paper_item"
	"sandbox/pkg/pagination"
)

// WorkPaperItemHandler handles HTTP requests for work paper items
type WorkPaperItemHandler struct {
	createUseCase *work_paper_item.CreateWorkPaperItemUseCase
	getUseCase    *work_paper_item.GetWorkPaperItemUseCase
	updateUseCase *work_paper_item.UpdateWorkPaperItemUseCase
	deleteUseCase *work_paper_item.DeleteWorkPaperItemUseCase
	listUseCase   *work_paper_item.ListWorkPaperItemsUseCase
	validator     *validator.Validate
}

// NewWorkPaperItemHandler creates a new handler instance
func NewWorkPaperItemHandler(
	createUseCase *work_paper_item.CreateWorkPaperItemUseCase,
	getUseCase *work_paper_item.GetWorkPaperItemUseCase,
	updateUseCase *work_paper_item.UpdateWorkPaperItemUseCase,
	deleteUseCase *work_paper_item.DeleteWorkPaperItemUseCase,
	listUseCase *work_paper_item.ListWorkPaperItemsUseCase,
) *WorkPaperItemHandler {
	return &WorkPaperItemHandler{
		createUseCase: createUseCase,
		getUseCase:    getUseCase,
		updateUseCase: updateUseCase,
		deleteUseCase: deleteUseCase,
		listUseCase:   listUseCase,
		validator:     validator.New(),
	}
}

// CreateWorkPaperItem creates a new work paper item
// @Summary Create Work Paper Item
// @Description Creates a new work paper item
// @Tags desk
// @Accept json
// @Produce json
// @Param request body work_paper_item.Request true "Create Work Paper Item Request"
// @Success 201 {object} StandardResponse{data=work_paper_item.Response}
// @Failure 400 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /api/v1/desk/work-paper-items [post]
func (h *WorkPaperItemHandler) CreateWorkPaperItem(c *fiber.Ctx) error {
	var req work_paper_item.Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
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
	response, err := h.createUseCase.Execute(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create work paper item",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetWorkPaperItem gets a work paper item by ID
// @Summary Get Work Paper Item
// @Description Gets a work paper item by its ID
// @Tags desk
// @Accept json
// @Produce json
// @Param id path string true "Work Paper Item ID"
// @Success 200 {object} StandardResponse{data=work_paper_item.GetResponse}
// @Failure 400 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /api/v1/desk/work-paper-items/{id} [get]
func (h *WorkPaperItemHandler) GetWorkPaperItem(c *fiber.Ctx) error {
	itemID := c.Params("id")
	if itemID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing work paper item ID",
			"details": "Work paper item ID is required",
		})
	}

	req := work_paper_item.GetRequest{
		ID: itemID,
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
	response, err := h.getUseCase.Execute(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get work paper item",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// ListWorkPaperItems lists work paper items
// @Summary List Work Paper Items
// @Description Lists work paper items with pagination and filtering
// @Tags desk
// @Accept json
// @Produce json
// @Param search query string false "Search term"
// @Param type query string false "Filter by type (A, B, C)"
// @Param is_active query bool false "Filter by active status"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Success 200 {object} pagination.PagedResponse
// @Failure 400 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /api/v1/desk/work-paper-items [get]
func (h *WorkPaperItemHandler) ListWorkPaperItems(c *fiber.Ctx) error {
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

	ctx := context.Background()
	workPaperItems, pagedResponse, err := h.listUseCase.Execute(ctx, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":     true,
		"data":        workPaperItems,
		"page":        pagedResponse.Page,
		"limit":       pagedResponse.Limit,
		"total_items": pagedResponse.TotalItems,
		"total_pages": pagedResponse.TotalPages,
	})
}

// UpdateWorkPaperItem updates an existing work paper item
// @Summary Update Work Paper Item
// @Description Updates an existing work paper item
// @Tags desk
// @Accept json
// @Produce json
// @Param id path string true "Work Paper Item ID"
// @Param request body work_paper_item.UpdateRequest true "Update Work Paper Item Request"
// @Success 200 {object} StandardResponse{data=work_paper_item.UpdateResponse}
// @Failure 400 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /api/v1/desk/work-paper-items/{id} [put]
func (h *WorkPaperItemHandler) UpdateWorkPaperItem(c *fiber.Ctx) error {
	itemID := c.Params("id")
	if itemID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing work paper item ID",
			"details": "Work paper item ID is required",
		})
	}

	var req work_paper_item.UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Set ID from path parameter
	req.ID = itemID

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	// Execute use case
	ctx := context.Background()
	response, err := h.updateUseCase.Execute(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update work paper item",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// DeleteWorkPaperItem deletes a work paper item
// @Summary Delete Work Paper Item
// @Description Deletes a work paper item (soft delete)
// @Tags desk
// @Accept json
// @Produce json
// @Param id path string true "Work Paper Item ID"
// @Success 200 {object} StandardResponse{data=work_paper_item.DeleteResponse}
// @Failure 400 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /api/v1/desk/work-paper-items/{id} [delete]
func (h *WorkPaperItemHandler) DeleteWorkPaperItem(c *fiber.Ctx) error {
	itemID := c.Params("id")
	if itemID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing work paper item ID",
			"details": "Work paper item ID is required",
		})
	}

	req := work_paper_item.DeleteRequest{
		ID: itemID,
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
	response, err := h.deleteUseCase.Execute(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete work paper item",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// Backward compatibility methods (deprecated)

// CreateMasterLakipItem creates a new master LAKIP item (deprecated)
// @Summary Create Master LAKIP Item (Deprecated)
// @Description Creates a new master LAKIP item (deprecated - use CreateWorkPaperItem instead)
// @Tags desk
// @Accept json
// @Produce json
// @Param request body work_paper_item.Request true "Create Master LAKIP Item Request"
// @Success 201 {object} StandardResponse{data=work_paper_item.Response}
// @Failure 400 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /api/v1/desk/master-lakip-items [post]
func (h *WorkPaperItemHandler) CreateMasterLakipItem(c *fiber.Ctx) error {
	return h.CreateWorkPaperItem(c)
}

// ListMasterLakipItems lists master LAKIP items (deprecated)
// @Summary List Master LAKIP Items (Deprecated)
// @Description Lists master LAKIP items with pagination and filtering (deprecated - use ListWorkPaperItems instead)
// @Tags desk
// @Accept json
// @Produce json
// @Param search query string false "Search term"
// @Param type query string false "Filter by type (A, B, C)"
// @Param is_active query bool false "Filter by active status"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} StandardResponse{data=work_paper_item.ListResponse}
// @Failure 400 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /api/v1/desk/master-lakip-items [get]
func (h *WorkPaperItemHandler) ListMasterLakipItems(c *fiber.Ctx) error {
	return h.ListWorkPaperItems(c)
}

// Backward compatibility factory function (deprecated)
func NewMasterLakipItemHandler(
	createUseCase *work_paper_item.CreateWorkPaperItemUseCase,
	getUseCase *work_paper_item.GetWorkPaperItemUseCase,
	updateUseCase *work_paper_item.UpdateWorkPaperItemUseCase,
	deleteUseCase *work_paper_item.DeleteWorkPaperItemUseCase,
	listUseCase *work_paper_item.ListWorkPaperItemsUseCase,
) *WorkPaperItemHandler {
	return NewWorkPaperItemHandler(createUseCase, getUseCase, updateUseCase, deleteUseCase, listUseCase)
}
