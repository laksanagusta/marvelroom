package desk

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"sandbox/internal/usecase/work_paper"
)

// WorkPaperHandler handles HTTP requests for work paper
type WorkPaperHandler struct {
	createUseCase          *work_paper.CreateWorkPaperUseCase
	checkDocumentUseCase   *work_paper.CheckWorkPaperNoteUseCase
	listUseCase            *work_paper.ListWorkPapersUseCase
	getDetailsUseCase      *work_paper.GetWorkPaperDetailsUseCase
	updateStatusUseCase     *work_paper.UpdateWorkPaperStatusUseCase
	updateWorkPaperNoteCase *work_paper.UpdateWorkPaperNoteUseCase
	validator              *validator.Validate
}

// NewWorkPaperHandler creates a new handler instance
func NewWorkPaperHandler(
	createUseCase *work_paper.CreateWorkPaperUseCase,
	checkDocumentUseCase *work_paper.CheckWorkPaperNoteUseCase,
	listUseCase *work_paper.ListWorkPapersUseCase,
	getDetailsUseCase *work_paper.GetWorkPaperDetailsUseCase,
	updateStatusUseCase *work_paper.UpdateWorkPaperStatusUseCase,
	updateWorkPaperNoteCase *work_paper.UpdateWorkPaperNoteUseCase,
) *WorkPaperHandler {
	return &WorkPaperHandler{
		createUseCase:          createUseCase,
		checkDocumentUseCase:   checkDocumentUseCase,
		listUseCase:            listUseCase,
		getDetailsUseCase:      getDetailsUseCase,
		updateStatusUseCase:     updateStatusUseCase,
		updateWorkPaperNoteCase: updateWorkPaperNoteCase,
		validator:              validator.New(),
	}
}

// CreateWorkPaper creates a new work paper
// @Summary Create Work Paper
// @Description Creates a new work paper for an organization and semester
// @Tags desk
// @Accept json
// @Produce json
// @Param request body work_paper.CreateRequest true "Create Work Paper Request"
// @Success 201 {object} StandardResponse{data=work_paper.CreateResponse}
// @Failure 400 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /api/v1/desk/work-papers [post]
func (h *WorkPaperHandler) CreateWorkPaper(c *fiber.Ctx) error {
	var req work_paper.CreateRequest
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
			"error":   "Failed to create work paper",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// ListWorkPapers lists work papers with pagination and filtering
// @Summary List Work Papers
// @Description Lists work papers with optional filters for organization, year, semester, and status
// @Tags desk
// @Accept json
// @Produce json
// @Param organization_id query string false "Organization ID filter"
// @Param year query int false "Year filter"
// @Param semester query int false "Semester filter (1 or 2)"
// @Param status query string false "Status filter"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} StandardResponse{data=work_paper.ListResponse}
// @Failure 400 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /api/v1/desk/work-papers [get]
func (h *WorkPaperHandler) ListWorkPapers(c *fiber.Ctx) error {
	// Parse query parameters
	req := work_paper.ListRequest{
		Page:     c.QueryInt("page", 1),
		PageSize: c.QueryInt("page_size", 10),
	}

	// Set optional filters
	if orgID := c.Query("organization_id"); orgID != "" {
		req.OrganizationID = orgID
	}
	if year := c.QueryInt("year", 0); year != 0 {
		req.Year = &year
	}
	if semester := c.QueryInt("semester", 0); semester != 0 {
		req.Semester = &semester
	}
	if status := c.Query("status"); status != "" {
		req.Status = status
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
	response, err := h.listUseCase.Execute(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to list work papers",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":     true,
		"data":        response.Data,
		"page":        response.Metadata.CurrentPage,
		"limit":       response.Metadata.PageSize,
		"total_items": response.Metadata.TotalCount,
		"total_pages": response.Metadata.TotalPage,
	})
}

// UpdateWorkPaperStatus updates the status of a work paper
// @Summary Update Work Paper Status
// @Description Updates the status of a work paper with validation for status transitions
// @Tags desk
// @Accept json
// @Produce json
// @Param id path string true "Work Paper ID"
// @Param request body work_paper.UpdateStatusRequest true "Update Status Request"
// @Success 200 {object} StandardResponse{data=work_paper.UpdateStatusResponse}
// @Failure 400 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /api/v1/desk/work-papers/{id}/status [put]
func (h *WorkPaperHandler) UpdateWorkPaperStatus(c *fiber.Ctx) error {
	// Get work paper ID from URL parameter
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Work Paper ID is required",
		})
	}

	// Parse request body
	var req work_paper.UpdateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Set ID from URL parameter
	req.ID = id

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	// Execute use case
	ctx := context.Background()
	response, err := h.updateStatusUseCase.Execute(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update work paper status",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetWorkPaperByID retrieves a work paper with all related details
// @Summary Get Work Paper by ID
// @Description Retrieves a work paper with its notes and signatures
// @Tags desk
// @Accept json
// @Produce json
// @Param id path string true "Work Paper ID"
// @Success 200 {object} StandardResponse{data=work_paper.GetWorkPaperDetailsResponse}
// @Failure 404 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /api/v1/desk/work-papers/{id} [get]
func (h *WorkPaperHandler) GetWorkPaperByID(c *fiber.Ctx) error {
	// Get work paper ID from URL parameter
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Work Paper ID is required",
		})
	}

	// Execute use case to get complete work paper details
	ctx := context.Background()
	response, err := h.getDetailsUseCase.Execute(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get work paper details",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetStatusTransitions returns allowed status transitions for a work paper
// @Summary Get Status Transitions
// @Description Returns the allowed status transitions for a given current status
// @Tags desk
// @Accept json
// @Produce json
// @Param current_status query string true "Current status"
// @Success 200 {object} StandardResponse{data=[]string}
// @Failure 400 {object} StandardResponse
// @Router /api/v1/desk/work-papers/status-transitions [get]
func (h *WorkPaperHandler) GetStatusTransitions(c *fiber.Ctx) error {
	currentStatus := c.Query("current_status")
	if currentStatus == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "current_status parameter is required",
		})
	}

	// Get allowed transitions
	transitions := work_paper.GetStatusTransitions(currentStatus)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    transitions,
	})
}

// CheckWorkPaperNote checks a work paper note using LLM
// @Summary Check Work Paper Note
// @Description Checks a work paper note document using LLM
// @Tags desk
// @Accept json
// @Produce json
// @Param request body work_paper.CheckRequest true "Check Work Paper Note Request"
// @Success 200 {object} StandardResponse{data=work_paper.CheckResponse}
// @Failure 400 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /api/v1/desk/work-paper-notes/check [post]
func (h *WorkPaperHandler) CheckWorkPaperNote(c *fiber.Ctx) error {
	var req work_paper.CheckRequest
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
	response, err := h.checkDocumentUseCase.Execute(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to check work paper note",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// UpdateWorkPaperNote updates a work paper note
// @Summary Update Work Paper Note
// @Description Updates a work paper note with Google Drive link, validation status, and/or notes
// @Tags desk
// @Accept json
// @Produce json
// @Param id path string true "Work Paper Note ID"
// @Param request body work_paper.UpdateWorkPaperNoteRequest true "Update Work Paper Note Request"
// @Success 200 {object} StandardResponse{data=work_paper.UpdateWorkPaperNoteResponse}
// @Failure 400 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /api/v1/desk/work-paper-notes/{id} [put]
func (h *WorkPaperHandler) UpdateWorkPaperNote(c *fiber.Ctx) error {
	// Get work paper note ID from URL parameter
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Work Paper Note ID is required",
		})
	}

	// Parse request body
	var req work_paper.UpdateWorkPaperNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Set ID from URL parameter
	req.ID = id

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	// Execute use case
	ctx := context.Background()
	response, err := h.updateWorkPaperNoteCase.Execute(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update work paper note",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// Backward compatibility methods (deprecated)

// CreatePaperWork creates a new paper work (deprecated)
// @Summary Create Paper Work (Deprecated)
// @Description Creates a new paper work for an organization and semester (deprecated - use CreateWorkPaper instead)
// @Tags desk
// @Accept json
// @Produce json
// @Param request body work_paper.CreateRequest true "Create Paper Work Request"
// @Success 201 {object} StandardResponse{data=work_paper.CreateResponse}
// @Failure 400 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /api/v1/desk/paper-works [post]
func (h *WorkPaperHandler) CreatePaperWork(c *fiber.Ctx) error {
	return h.CreateWorkPaper(c)
}

// CheckDocument checks a document using LLM (deprecated)
// @Summary Check Document (Deprecated)
// @Description Checks a document using LLM (deprecated - use CheckWorkPaperNote instead)
// @Tags desk
// @Accept json
// @Produce json
// @Param request body work_paper.CheckRequest true "Check Document Request"
// @Success 200 {object} StandardResponse{data=work_paper.CheckResponse}
// @Failure 400 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /api/v1/desk/paper-work-items/check [post]
func (h *WorkPaperHandler) CheckDocument(c *fiber.Ctx) error {
	// For backward compatibility, convert the old request format to new format
	var req struct {
		ItemID string `json:"item_id" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Convert to new request format
	newReq := work_paper.CheckRequest{
		NoteID: req.ItemID,
	}

	// Validate request
	if err := h.validator.Struct(&newReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	// Execute use case
	ctx := context.Background()
	response, err := h.checkDocumentUseCase.Execute(ctx, newReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to check document",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// Backward compatibility factory function (deprecated)
func NewPaperWorkHandler(
	createUseCase *work_paper.CreateWorkPaperUseCase,
	checkDocumentUseCase *work_paper.CheckWorkPaperNoteUseCase,
) *WorkPaperHandler {
	return NewWorkPaperHandler(createUseCase, checkDocumentUseCase, nil, nil, nil, nil)
}