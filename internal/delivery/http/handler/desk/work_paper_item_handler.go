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
	listUseCase   *work_paper_item.ListWorkPaperItemsUseCase
	validator     *validator.Validate
}

// NewWorkPaperItemHandler creates a new handler instance
func NewWorkPaperItemHandler(
	createUseCase *work_paper_item.CreateWorkPaperItemUseCase,
	listUseCase *work_paper_item.ListWorkPaperItemsUseCase,
) *WorkPaperItemHandler {
	return &WorkPaperItemHandler{
		createUseCase: createUseCase,
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
	listUseCase *work_paper_item.ListWorkPaperItemsUseCase,
) *WorkPaperItemHandler {
	return NewWorkPaperItemHandler(createUseCase, listUseCase)
}