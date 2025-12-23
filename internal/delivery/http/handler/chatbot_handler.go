package handler

import (
	"io"

	"sandbox/internal/domain/entity"
	chatbotUC "sandbox/internal/usecase/chatbot"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ChatbotHandler handles chatbot HTTP requests
type ChatbotHandler struct {
	createKnowledgeBaseUseCase *chatbotUC.CreateKnowledgeBaseUseCase
	listKnowledgeBasesUseCase  *chatbotUC.ListKnowledgeBasesUseCase
	getKnowledgeBaseUseCase    *chatbotUC.GetKnowledgeBaseUseCase
	deleteKnowledgeBaseUseCase *chatbotUC.DeleteKnowledgeBaseUseCase
	uploadFilesUseCase         *chatbotUC.UploadFilesUseCase
	syncDocumentStatusUseCase  *chatbotUC.SyncDocumentStatusUseCase
	createChatSessionUseCase   *chatbotUC.CreateChatSessionUseCase
	listChatSessionsUseCase    *chatbotUC.ListChatSessionsUseCase
	deleteChatSessionUseCase   *chatbotUC.DeleteChatSessionUseCase
	sendMessageUseCase         *chatbotUC.SendMessageUseCase
	getChatHistoryUseCase      *chatbotUC.GetChatHistoryUseCase
	validation                 *validator.Validate
}

// NewChatbotHandler creates a new ChatbotHandler
func NewChatbotHandler(
	createKnowledgeBaseUseCase *chatbotUC.CreateKnowledgeBaseUseCase,
	listKnowledgeBasesUseCase *chatbotUC.ListKnowledgeBasesUseCase,
	getKnowledgeBaseUseCase *chatbotUC.GetKnowledgeBaseUseCase,
	deleteKnowledgeBaseUseCase *chatbotUC.DeleteKnowledgeBaseUseCase,
	uploadFilesUseCase *chatbotUC.UploadFilesUseCase,
	syncDocumentStatusUseCase *chatbotUC.SyncDocumentStatusUseCase,
	createChatSessionUseCase *chatbotUC.CreateChatSessionUseCase,
	listChatSessionsUseCase *chatbotUC.ListChatSessionsUseCase,
	deleteChatSessionUseCase *chatbotUC.DeleteChatSessionUseCase,
	sendMessageUseCase *chatbotUC.SendMessageUseCase,
	getChatHistoryUseCase *chatbotUC.GetChatHistoryUseCase,
) *ChatbotHandler {
	return &ChatbotHandler{
		createKnowledgeBaseUseCase: createKnowledgeBaseUseCase,
		listKnowledgeBasesUseCase:  listKnowledgeBasesUseCase,
		getKnowledgeBaseUseCase:    getKnowledgeBaseUseCase,
		deleteKnowledgeBaseUseCase: deleteKnowledgeBaseUseCase,
		uploadFilesUseCase:         uploadFilesUseCase,
		syncDocumentStatusUseCase:  syncDocumentStatusUseCase,
		createChatSessionUseCase:   createChatSessionUseCase,
		listChatSessionsUseCase:    listChatSessionsUseCase,
		deleteChatSessionUseCase:   deleteChatSessionUseCase,
		sendMessageUseCase:         sendMessageUseCase,
		getChatHistoryUseCase:      getChatHistoryUseCase,
		validation:                 validator.New(),
	}
}

// CreateKnowledgeBaseRequest is the request body
type CreateKnowledgeBaseRequest struct {
	Name string `json:"name" validate:"required"`
}

// CreateKnowledgeBase handles POST /api/v1/chatbot/knowledge-bases
func (h *ChatbotHandler) CreateKnowledgeBase(c *fiber.Ctx) error {
	user := c.Locals("authenticatedUser").(*entity.AuthenticatedUser)

	var req CreateKnowledgeBaseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validation.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	output, err := h.createKnowledgeBaseUseCase.Execute(c.Context(), chatbotUC.CreateKnowledgeBaseInput{
		Name: req.Name,
	}, *user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": output.KnowledgeBase,
	})
}

// ListKnowledgeBases handles GET /api/v1/chatbot/knowledge-bases
func (h *ChatbotHandler) ListKnowledgeBases(c *fiber.Ctx) error {
	user := c.Locals("authenticatedUser").(*entity.AuthenticatedUser)

	output, err := h.listKnowledgeBasesUseCase.Execute(c.Context(), *user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": output.KnowledgeBases,
	})
}

// GetKnowledgeBase handles GET /api/v1/chatbot/knowledge-bases/:id
func (h *ChatbotHandler) GetKnowledgeBase(c *fiber.Ctx) error {
	user := c.Locals("authenticatedUser").(*entity.AuthenticatedUser)
	kbID := c.Params("id")

	output, err := h.getKnowledgeBaseUseCase.Execute(c.Context(), chatbotUC.GetKnowledgeBaseInput{
		KnowledgeBaseID: kbID,
	}, *user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"knowledge_base": output.KnowledgeBase,
			"documents":      output.Documents,
		},
	})
}

// DeleteKnowledgeBase handles DELETE /api/v1/chatbot/knowledge-bases/:id
func (h *ChatbotHandler) DeleteKnowledgeBase(c *fiber.Ctx) error {
	user := c.Locals("authenticatedUser").(*entity.AuthenticatedUser)
	kbID := c.Params("id")

	err := h.deleteKnowledgeBaseUseCase.Execute(c.Context(), chatbotUC.DeleteKnowledgeBaseInput{
		KnowledgeBaseID: kbID,
	}, *user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// UploadFiles handles POST /api/v1/chatbot/knowledge-bases/:id/files
func (h *ChatbotHandler) UploadFiles(c *fiber.Ctx) error {
	user := c.Locals("authenticatedUser").(*entity.AuthenticatedUser)
	kbID := c.Params("id")

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid multipart form",
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No files provided",
		})
	}

	var uploadFiles []chatbotUC.UploadFileInput
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to open file: " + fileHeader.Filename,
			})
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to read file: " + fileHeader.Filename,
			})
		}

		uploadFiles = append(uploadFiles, chatbotUC.UploadFileInput{
			FileName: fileHeader.Filename,
			Content:  content,
			MimeType: fileHeader.Header.Get("Content-Type"),
		})
	}

	output, err := h.uploadFilesUseCase.Execute(c.Context(), chatbotUC.UploadFilesInput{
		KnowledgeBaseID: kbID,
		Files:           uploadFiles,
	}, *user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Format errors as strings
	var errorStrings []string
	for _, e := range output.Errors {
		errorStrings = append(errorStrings, e.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"documents": output.Documents,
			"errors":    errorStrings,
		},
	})
}

// SyncDocumentStatus handles POST /api/v1/chatbot/knowledge-bases/:id/sync-status
func (h *ChatbotHandler) SyncDocumentStatus(c *fiber.Ctx) error {
	user := c.Locals("authenticatedUser").(*entity.AuthenticatedUser)
	kbID := c.Params("id")

	output, err := h.syncDocumentStatusUseCase.Execute(c.Context(), chatbotUC.SyncDocumentStatusInput{
		KnowledgeBaseID: kbID,
	}, *user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"documents":       output.Documents,
			"updated_count":   output.UpdatedCount,
			"processing_docs": output.ProcessingDocs,
			"active_docs":     output.ActiveDocs,
			"failed_docs":     output.FailedDocs,
		},
	})
}

// CreateChatSessionRequest is the request body
type CreateChatSessionRequest struct {
	KnowledgeBaseID string `json:"knowledge_base_id" validate:"required"`
	Title           string `json:"title"`
}

// CreateChatSession handles POST /api/v1/chatbot/sessions
func (h *ChatbotHandler) CreateChatSession(c *fiber.Ctx) error {
	user := c.Locals("authenticatedUser").(*entity.AuthenticatedUser)

	var req CreateChatSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validation.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	output, err := h.createChatSessionUseCase.Execute(c.Context(), chatbotUC.CreateChatSessionInput{
		KnowledgeBaseID: req.KnowledgeBaseID,
		Title:           req.Title,
	}, *user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": output.ChatSession,
	})
}

// ListChatSessions handles GET /api/v1/chatbot/sessions
func (h *ChatbotHandler) ListChatSessions(c *fiber.Ctx) error {
	user := c.Locals("authenticatedUser").(*entity.AuthenticatedUser)

	output, err := h.listChatSessionsUseCase.Execute(c.Context(), *user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": output.Sessions,
	})
}

// SendMessageRequest is the request body
type SendMessageRequest struct {
	Content string `json:"content" validate:"required"`
}

// SendMessage handles POST /api/v1/chatbot/sessions/:id/messages
func (h *ChatbotHandler) SendMessage(c *fiber.Ctx) error {
	user := c.Locals("authenticatedUser").(*entity.AuthenticatedUser)
	sessionID := c.Params("id")

	var req SendMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validation.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	output, err := h.sendMessageUseCase.Execute(c.Context(), chatbotUC.SendMessageInput{
		ChatSessionID: sessionID,
		Content:       req.Content,
	}, *user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"user_message":      output.UserMessage,
			"assistant_message": output.AssistantMessage,
		},
	})
}

// GetChatHistory handles GET /api/v1/chatbot/sessions/:id/messages
func (h *ChatbotHandler) GetChatHistory(c *fiber.Ctx) error {
	user := c.Locals("authenticatedUser").(*entity.AuthenticatedUser)
	sessionID := c.Params("id")

	output, err := h.getChatHistoryUseCase.Execute(c.Context(), chatbotUC.GetChatHistoryInput{
		ChatSessionID: sessionID,
	}, *user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": output.Messages,
	})
}

// DeleteChatSession handles DELETE /api/v1/chatbot/sessions/:id
func (h *ChatbotHandler) DeleteChatSession(c *fiber.Ctx) error {
	user := c.Locals("authenticatedUser").(*entity.AuthenticatedUser)
	sessionID := c.Params("id")

	err := h.deleteChatSessionUseCase.Execute(c.Context(), chatbotUC.DeleteChatSessionInput{
		SessionID: sessionID,
	}, *user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
