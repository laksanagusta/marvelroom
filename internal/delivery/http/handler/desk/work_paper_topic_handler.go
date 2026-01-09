package desk

import (
	"io"
	"strconv"
	"strings"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/service"
	"sandbox/pkg/pagination"

	"github.com/gofiber/fiber/v2"
)

// WorkPaperTopicHandler handles HTTP requests for work paper topics
type WorkPaperTopicHandler struct {
	deskService service.DeskService
}

// NewWorkPaperTopicHandler creates a new handler instance
func NewWorkPaperTopicHandler(deskService service.DeskService) *WorkPaperTopicHandler {
	return &WorkPaperTopicHandler{
		deskService: deskService,
	}
}

// CreateTopicRequest represents the request for creating a topic
type CreateTopicRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateTopicRequest represents the request for updating a topic
type UpdateTopicRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

// TopicResponse represents the response for a topic
type TopicResponse struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Description     *string `json:"description,omitempty"`
	TemplatePath    *string `json:"template_path,omitempty"`
	TemplateVersion int     `json:"template_version"`
	IsActive        bool    `json:"is_active"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// toTopicResponse converts entity to response DTO
func toTopicResponse(topic *entity.WorkPaperTopic) *TopicResponse {
	return &TopicResponse{
		ID:              topic.ID.String(),
		Name:            topic.Name,
		Description:     topic.Description,
		TemplatePath:    topic.TemplatePath,
		TemplateVersion: topic.TemplateVersion,
		IsActive:        topic.IsActive,
		CreatedAt:       topic.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       topic.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// CreateWorkPaperTopic creates a new work paper topic
func (h *WorkPaperTopicHandler) CreateWorkPaperTopic(c *fiber.Ctx) error {
	var req CreateTopicRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name is required"})
	}

	serviceReq := &service.CreateWorkPaperTopicRequest{
		Name:        req.Name,
		Description: req.Description,
	}

	topic, err := h.deskService.CreateWorkPaperTopic(c.Context(), serviceReq)
	if err != nil {
		if err == entity.ErrDuplicateWorkPaperTopicName {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "topic with this name already exists"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": toTopicResponse(topic)})
}

// GetWorkPaperTopic gets a work paper topic by ID
func (h *WorkPaperTopicHandler) GetWorkPaperTopic(c *fiber.Ctx) error {
	id := c.Params("id")

	topic, err := h.deskService.GetWorkPaperTopic(c.Context(), id)
	if err != nil {
		if err == entity.ErrWorkPaperTopicNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "topic not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": toTopicResponse(topic)})
}

// UpdateWorkPaperTopic updates a work paper topic
func (h *WorkPaperTopicHandler) UpdateWorkPaperTopic(c *fiber.Ctx) error {
	id := c.Params("id")

	var req UpdateTopicRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name is required"})
	}

	serviceReq := &service.UpdateWorkPaperTopicRequest{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	topic, err := h.deskService.UpdateWorkPaperTopic(c.Context(), serviceReq)
	if err != nil {
		if err == entity.ErrWorkPaperTopicNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "topic not found"})
		}
		if err == entity.ErrDuplicateWorkPaperTopicName {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "topic with this name already exists"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": toTopicResponse(topic)})
}

// DeleteWorkPaperTopic deletes a work paper topic (soft delete)
func (h *WorkPaperTopicHandler) DeleteWorkPaperTopic(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.deskService.DeleteWorkPaperTopic(c.Context(), id)
	if err != nil {
		if err == entity.ErrWorkPaperTopicNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "topic not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ListWorkPaperTopics lists all work paper topics with pagination
func (h *WorkPaperTopicHandler) ListWorkPaperTopics(c *fiber.Ctx) error {
	// Parse pagination params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	params := &pagination.QueryParams{
		Pagination: pagination.Pagination{
			Page:  page,
			Limit: limit,
		},
		Filters: []pagination.Filter{},
	}

	// Add is_active filter if provided
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		params.Filters = append(params.Filters, pagination.Filter{
			Field:    "is_active",
			Operator: "eq",
			Value:    isActive,
		})
	}

	topics, total, err := h.deskService.ListWorkPaperTopics(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Convert to response
	var responses []*TopicResponse
	for _, topic := range topics {
		responses = append(responses, toTopicResponse(topic))
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return c.JSON(fiber.Map{
		"data": responses,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total_items": total,
			"total_pages": totalPages,
		},
	})
}

// GetActiveWorkPaperTopics gets all active work paper topics
func (h *WorkPaperTopicHandler) GetActiveWorkPaperTopics(c *fiber.Ctx) error {
	topics, err := h.deskService.GetActiveWorkPaperTopics(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Convert to response
	var responses []*TopicResponse
	for _, topic := range topics {
		responses = append(responses, toTopicResponse(topic))
	}

	return c.JSON(fiber.Map{"data": responses})
}

// UploadTemplate handles uploading an Excel template for a work paper topic
func (h *WorkPaperTopicHandler) UploadTemplate(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get the uploaded file
	file, err := c.FormFile("template")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "template file is required"})
	}

	// Open the file
	fileHeader, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to open uploaded file"})
	}
	defer fileHeader.Close()

	// Read file content
	fileContent, err := io.ReadAll(fileHeader)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to read file content"})
	}

	// Get content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		// Try to detect from filename extension
		if strings.HasSuffix(strings.ToLower(file.Filename), ".xlsx") {
			contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		} else if strings.HasSuffix(strings.ToLower(file.Filename), ".xls") {
			contentType = "application/vnd.ms-excel"
		}
	}

	serviceReq := &service.UploadWorkPaperTopicTemplateRequest{
		TopicID:     id,
		FileName:    file.Filename,
		FileContent: fileContent,
		ContentType: contentType,
	}

	topic, err := h.deskService.UploadWorkPaperTopicTemplate(c.Context(), serviceReq)
	if err != nil {
		if err == entity.ErrWorkPaperTopicNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "topic not found"})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Template uploaded successfully",
		"data":    toTopicResponse(topic),
	})
}

// DeleteTemplate handles deleting the Excel template from a work paper topic
func (h *WorkPaperTopicHandler) DeleteTemplate(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.deskService.DeleteWorkPaperTopicTemplate(c.Context(), id)
	if err != nil {
		if err == entity.ErrWorkPaperTopicNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "topic not found"})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Template deleted successfully",
	})
}

// DownloadTemplate handles downloading the Excel template from a work paper topic
func (h *WorkPaperTopicHandler) DownloadTemplate(c *fiber.Ctx) error {
	id := c.Params("id")

	topic, err := h.deskService.GetWorkPaperTopic(c.Context(), id)
	if err != nil {
		if err == entity.ErrWorkPaperTopicNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "topic not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if topic.TemplatePath == nil || *topic.TemplatePath == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "no template file found for this topic"})
	}

	// Send the file
	return c.SendFile(*topic.TemplatePath)
}
