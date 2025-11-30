package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"sandbox/internal/delivery/http/middleware"
	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/service"
	workPaperSignatureUC "sandbox/internal/usecase/work_paper_signature"
	"sandbox/pkg/pagination"
)

type WorkPaperSignatureHandler struct {
	deskService                                service.DeskService
	listWorkPaperSignaturesUseCase             *workPaperSignatureUC.ListWorkPaperSignaturesUseCase
	getWorkPaperSignaturesByWorkPaperIDUseCase *workPaperSignatureUC.GetWorkPaperSignaturesByWorkPaperIDUseCase
	createDigitalSignatureUseCase              *workPaperSignatureUC.CreateDigitalSignatureUseCase
	verifyDigitalSignatureUseCase              *workPaperSignatureUC.VerifyDigitalSignatureUseCase
	validation                                 *validator.Validate
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ListWorkPaperSignaturesResponse represents response for listing work paper signatures with pagination
type ListWorkPaperSignaturesResponse struct {
	Signatures  []*entity.WorkPaperSignature `json:"data"`
	TotalItems  int64                        `json:"total_items"`
	TotalPages  int                          `json:"total_pages"`
	CurrentPage int                          `json:"current_page"`
	Limit       int                          `json:"limit"`
}

// CreateWorkPaperSignature creates a new work paper signature
// @Summary Create Work Paper Signature
// @Description Creates a new signature request for a work paper
// @Tags work-paper-signatures
// @Accept json
// @Produce json
// @Param request body service.CreateWorkPaperSignatureRequest true "Signature request"
// @Success 200 {object} entity.WorkPaperSignature
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/work-paper-signatures [post]
func (h *WorkPaperSignatureHandler) CreateWorkPaperSignature(c *fiber.Ctx) error {
	var req service.CreateWorkPaperSignatureRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	ctx := context.Background()
	signature, err := h.deskService.CreateWorkPaperSignature(ctx, &req)
	if err != nil {
		switch err {
		case entity.ErrWorkPaperNoteNotFound:
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "Work paper not found",
				Message: err.Error(),
				Code:    fiber.StatusNotFound,
			})
		case entity.ErrDuplicateSignature:
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
				Error:   "Duplicate signature",
				Message: "A signature already exists for this user and work paper",
				Code:    fiber.StatusConflict,
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to create signature",
				Message: err.Error(),
				Code:    fiber.StatusInternalServerError,
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(SuccessResponse{
		Success: true,
		Message: "Work paper signature created successfully",
		Data:    signature,
	})
}

// GetWorkPaperSignature gets a work paper signature by ID
// @Summary Get Work Paper Signature
// @Description Gets a work paper signature by ID
// @Tags work-paper-signatures
// @Produce json
// @Param id path string true "Signature ID"
// @Success 200 {object} entity.WorkPaperSignature
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/work-paper-signatures/{id} [get]
func (h *WorkPaperSignatureHandler) GetWorkPaperSignature(c *fiber.Ctx) error {
	signatureID := c.Params("id")
	if signatureID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Missing signature ID",
			Message: "Signature ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	ctx := context.Background()
	signature, err := h.deskService.GetWorkPaperSignature(ctx, signatureID)
	if err != nil {
		switch err {
		case entity.ErrSignatureNotFound:
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "Signature not found",
				Message: err.Error(),
				Code:    fiber.StatusNotFound,
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to get signature",
				Message: err.Error(),
				Code:    fiber.StatusInternalServerError,
			})
		}
	}

	return c.JSON(SuccessResponse{
		Success: true,
		Message: "Work paper signature retrieved successfully",
		Data:    signature,
	})
}

// ListWorkPapersWithSignatures lists all work papers with their signatures
// @Summary List Work Papers with Signatures
// @Description Gets all work papers with their signature details
// @Tags work-paper-signatures
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page"
// @Param status query string false "Filter by work paper status"
// @Param organizationId query string false "Filter by organization ID"
// @Success 200 {object} SuccessResponse{data=[]service.WorkPaperWithSignatures}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/work-paper-signatures/work-papers [get]
func (h *WorkPaperSignatureHandler) ListWorkPapersWithSignatures(c *fiber.Ctx) error {
	// Parse query parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	status := c.Query("status")
	organizationID := c.Query("organizationId")

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Create context with timeout
	ctx := context.Background()

	// Get work papers with their signatures
	workPapers, err := h.deskService.GetWorkPapersWithSignatures(ctx, page, limit, status, organizationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to retrieve work papers with signatures",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	return c.JSON(SuccessResponse{
		Success: true,
		Message: "Work papers with signatures retrieved successfully",
		Data:    workPapers,
	})
}

// SignWorkPaper signs a work paper
// @Summary Sign Work Paper
// @Description Signs a work paper using authenticated user credentials and generates SHA256 hash
// @Tags work-paper-signatures
// @Produce json
// @Param id path string true "Signature ID"
// @Success 200 {object} entity.WorkPaperSignature
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/work-paper-signatures/{id}/sign [post]
func (h *WorkPaperSignatureHandler) SignWorkPaper(c *fiber.Ctx) error {
	signatureID := c.Params("id")
	if signatureID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Missing signature ID",
			Message: "Signature ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Get authenticated user from middleware context
	user, err := middleware.GetAuthenticatedUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Authentication required",
			Message: "User authentication is required to sign work paper",
			Code:    fiber.StatusUnauthorized,
		})
	}

	ctx := context.Background()
	signature, err := h.deskService.SignWorkPaperWithUser(ctx, signatureID, user.ID)
	if err != nil {
		switch err {
		case entity.ErrSignatureNotFound:
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "Signature not found",
				Message: err.Error(),
				Code:    fiber.StatusNotFound,
			})
		case entity.ErrAlreadySigned:
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
				Error:   "Already signed",
				Message: "This signature has already been signed",
				Code:    fiber.StatusConflict,
			})
		case entity.ErrSignatureRejected:
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
				Error:   "Signature rejected",
				Message: "This signature has been rejected and cannot be signed",
				Code:    fiber.StatusConflict,
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to sign work paper",
				Message: err.Error(),
				Code:    fiber.StatusInternalServerError,
			})
		}
	}

	return c.JSON(SuccessResponse{
		Success: true,
		Message: "Work paper signed successfully",
		Data:    signature,
	})
}

// RejectWorkPaperSignature rejects a work paper signature
// @Summary Reject Work Paper Signature
// @Description Rejects a work paper signature
// @Tags work-paper-signatures
// @Accept json
// @Produce json
// @Param id path string true "Signature ID"
// @Param request body service.RejectWorkPaperSignatureRequest true "Reject request"
// @Success 200 {object} entity.WorkPaperSignature
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/work-paper-signatures/{id}/reject [post]
func (h *WorkPaperSignatureHandler) RejectWorkPaperSignature(c *fiber.Ctx) error {
	signatureID := c.Params("id")
	if signatureID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Missing signature ID",
			Message: "Signature ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	var req service.RejectWorkPaperSignatureRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	// Validate request
	if err := h.validation.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	ctx := context.Background()
	signature, err := h.deskService.RejectWorkPaperSignature(ctx, signatureID, &req)
	if err != nil {
		switch err {
		case entity.ErrSignatureNotFound:
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "Signature not found",
				Message: err.Error(),
				Code:    fiber.StatusNotFound,
			})
		case entity.ErrAlreadySigned:
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
				Error:   "Already signed",
				Message: "This signature has already been signed and cannot be rejected",
				Code:    fiber.StatusConflict,
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to reject signature",
				Message: err.Error(),
				Code:    fiber.StatusInternalServerError,
			})
		}
	}

	return c.JSON(SuccessResponse{
		Success: true,
		Message: "Work paper signature rejected successfully",
		Data:    signature,
	})
}

// ResetWorkPaperSignature resets a work paper signature to pending status
// @Summary Reset Work Paper Signature
// @Description Resets a work paper signature to pending status
// @Tags work-paper-signatures
// @Produce json
// @Param id path string true "Signature ID"
// @Success 200 {object} entity.WorkPaperSignature
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/work-paper-signatures/{id}/reset [post]
func (h *WorkPaperSignatureHandler) ResetWorkPaperSignature(c *fiber.Ctx) error {
	signatureID := c.Params("id")
	if signatureID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Missing signature ID",
			Message: "Signature ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	ctx := context.Background()
	signature, err := h.deskService.ResetWorkPaperSignature(ctx, signatureID)
	if err != nil {
		switch err {
		case entity.ErrSignatureNotFound:
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "Signature not found",
				Message: err.Error(),
				Code:    fiber.StatusNotFound,
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to reset signature",
				Message: err.Error(),
				Code:    fiber.StatusInternalServerError,
			})
		}
	}

	return c.JSON(SuccessResponse{
		Success: true,
		Message: "Work paper signature reset successfully",
		Data:    signature,
	})
}

// GetWorkPaperSignaturesByUserID gets all signatures for a specific user
// @Summary Get User's Work Paper Signatures
// @Description Gets all work paper signatures for a specific user
// @Tags work-paper-signatures
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {array} entity.WorkPaperSignature
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/{userId}/work-paper-signatures [get]
func (h *WorkPaperSignatureHandler) GetWorkPaperSignaturesByUserID(c *fiber.Ctx) error {
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Missing user ID",
			Message: "User ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	ctx := context.Background()
	signatures, err := h.deskService.GetWorkPaperSignaturesByUserID(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to get signatures by user ID",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	return c.JSON(SuccessResponse{
		Success: true,
		Message: "User's work paper signatures retrieved successfully",
		Data:    signatures,
	})
}

// ListWorkPaperSignatures lists work paper signatures with pagination and filtering
// @Summary List Work Paper Signatures
// @Description Lists work paper signatures with pagination and filtering (same as business trip list implementation)
// @Tags work-paper-signatures
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page"
// @Param user_id query string false "Filter by user ID"
// @Param status query string false "Filter by signature status (pending, signed, rejected)"
// @Param work_paper_id query string false "Filter by work paper ID"
// @Param sort_by query string false "Sort by field" default("created_at")
// @Param sort_dir query string false "Sort direction (asc, desc)" default("desc")
// @Success 200 {object} SuccessResponse{data=[]entity.WorkPaperSignature}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/work-paper-signatures [get]
func (h *WorkPaperSignatureHandler) ListWorkPaperSignatures(c *fiber.Ctx) error {
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

	workPaperSignatures, pagination, err := h.listWorkPaperSignaturesUseCase.Execute(context.Background(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	pagination.Data = workPaperSignatures

	return c.JSON(pagination)
}

// GetWorkPaperSignaturesByWorkPaperID gets all signatures for a specific work paper
// @Summary Get Work Paper Signatures by Work Paper ID
// @Description Gets all signatures for a specific work paper by its ID
// @Tags work-paper-signatures
// @Produce json
// @Param workPaperId path string true "Work Paper ID"
// @Success 200 {object} SuccessResponse{data=[]workPaperSignatureUC.WorkPaperSignatureResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/work-papers/{workPaperId}/signatures [get]
func (h *WorkPaperSignatureHandler) GetWorkPaperSignaturesByWorkPaperID(c *fiber.Ctx) error {
	workPaperID := c.Params("workPaperId")
	if workPaperID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Missing work paper ID",
			Message: "Work paper ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	ctx := context.Background()
	signatures, err := h.getWorkPaperSignaturesByWorkPaperIDUseCase.Execute(ctx, workPaperID)
	if err != nil {
		switch err {
		case workPaperSignatureUC.ErrInvalidWorkPaperID:
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "Invalid work paper ID",
				Message: err.Error(),
				Code:    fiber.StatusBadRequest,
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to get signatures by work paper ID",
				Message: err.Error(),
				Code:    fiber.StatusInternalServerError,
			})
		}
	}

	return c.JSON(SuccessResponse{
		Success: true,
		Message: "Work paper signatures retrieved successfully",
		Data:    signatures,
	})
}

// CreateDigitalSignature creates a digital signature for a work paper signature
// @Summary Create Digital Signature
// @Description Creates a certificate-based digital signature for a work paper signature
// @Tags work-paper-signatures
// @Produce json
// @Param id path string true "Signature ID"
// @Success 200 {object} workPaperSignatureUC.CreateDigitalSignatureResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/work-paper-signatures/{id}/digital-sign [post]
func (h *WorkPaperSignatureHandler) CreateDigitalSignature(c *fiber.Ctx) error {
	signatureID := c.Params("id")
	if signatureID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Missing signature ID",
			Message: "Signature ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Get authenticated user from context
	user, err := middleware.GetAuthenticatedUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
			Code:    fiber.StatusUnauthorized,
		})
	}

	// Create request with signature ID from URL and user ID from context
	req := &workPaperSignatureUC.CreateDigitalSignatureRequest{
		WorkPaperSignatureID: signatureID,
		UserID:               user.ID,
	}

	response, err := h.createDigitalSignatureUseCase.Execute(req)
	if err != nil {
		switch err {
		case workPaperSignatureUC.ErrWorkPaperSignatureNotFound:
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "Work paper signature not found",
				Message: err.Error(),
				Code:    fiber.StatusNotFound,
			})
		case workPaperSignatureUC.ErrCannotSignRejectedSignature:
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "Cannot sign rejected signature",
				Message: err.Error(),
				Code:    fiber.StatusBadRequest,
			})
		case entity.ErrAlreadySigned:
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
				Error:   "Already signed",
				Message: err.Error(),
				Code:    fiber.StatusConflict,
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to create digital signature",
				Message: err.Error(),
				Code:    fiber.StatusInternalServerError,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse{
		Success: true,
		Message: "Digital signature created successfully",
		Data:    response,
	})
}

// VerifyDigitalSignature verifies a digital signature
// @Summary Verify Digital Signature
// @Description Verifies a certificate-based digital signature for a work paper signature
// @Tags work-paper-signatures
// @Produce json
// @Param id path string true "Signature ID"
// @Success 200 {object} workPaperSignatureUC.VerifyDigitalSignatureResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/work-paper-signatures/{id}/verify [post]
func (h *WorkPaperSignatureHandler) VerifyDigitalSignature(c *fiber.Ctx) error {
	signatureID := c.Params("id")
	if signatureID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Missing signature ID",
			Message: "Signature ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Create request with signature ID from URL parameters
	req := &workPaperSignatureUC.VerifyDigitalSignatureRequest{
		WorkPaperSignatureID: signatureID,
	}

	response, err := h.verifyDigitalSignatureUseCase.Execute(req)
	if err != nil {
		switch err {
		case workPaperSignatureUC.ErrSignatureNotFound:
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "Signature not found",
				Message: err.Error(),
				Code:    fiber.StatusNotFound,
			})
		case workPaperSignatureUC.ErrNoDigitalSignature:
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "No digital signature found",
				Message: err.Error(),
				Code:    fiber.StatusNotFound,
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to verify digital signature",
				Message: err.Error(),
				Code:    fiber.StatusInternalServerError,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse{
		Success: true,
		Message: "Digital signature verification completed",
		Data:    response,
	})
}

// Updated constructor
func NewWorkPaperSignatureHandler(
	deskService service.DeskService,
	listWorkPaperSignaturesUseCase *workPaperSignatureUC.ListWorkPaperSignaturesUseCase,
	getWorkPaperSignaturesByWorkPaperIDUseCase *workPaperSignatureUC.GetWorkPaperSignaturesByWorkPaperIDUseCase,
	createDigitalSignatureUseCase *workPaperSignatureUC.CreateDigitalSignatureUseCase,
	verifyDigitalSignatureUseCase *workPaperSignatureUC.VerifyDigitalSignatureUseCase,
) *WorkPaperSignatureHandler {
	return &WorkPaperSignatureHandler{
		deskService:                                deskService,
		listWorkPaperSignaturesUseCase:             listWorkPaperSignaturesUseCase,
		getWorkPaperSignaturesByWorkPaperIDUseCase: getWorkPaperSignaturesByWorkPaperIDUseCase,
		createDigitalSignatureUseCase:              createDigitalSignatureUseCase,
		verifyDigitalSignatureUseCase:              verifyDigitalSignatureUseCase,
		validation:                                 validator.New(),
	}
}
