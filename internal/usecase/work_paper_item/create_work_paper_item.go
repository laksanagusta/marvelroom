package work_paper_item

import (
	"context"

	"github.com/google/uuid"

	"sandbox/internal/domain/service"
)

// parseUUID helper function to parse UUID string
func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// CreateWorkPaperItemUseCase handles the creation of work paper items
type CreateWorkPaperItemUseCase struct {
	deskService service.DeskService
}

// NewCreateWorkPaperItemUseCase creates a new use case instance
func NewCreateWorkPaperItemUseCase(deskService service.DeskService) *CreateWorkPaperItemUseCase {
	return &CreateWorkPaperItemUseCase{
		deskService: deskService,
	}
}

// Request represents the request payload for creating a work paper item
type Request struct {
	Type               string  `json:"type" validate:"required"`
	Number             string  `json:"number" validate:"required"`
	TopicID            string  `json:"topic_id"`
	DeskInstruction    string  `json:"desk_instruction" validate:"required"`
	ExpectedFolderName *string `json:"expected_folder_name,omitempty"`
	ParentID           string  `json:"parent_id,omitempty"`
	Level              int     `json:"level"`
}

// Response represents the response payload for creating a work paper item
type Response struct {
	ID                 string  `json:"id"`
	Type               string  `json:"type"`
	Number             string  `json:"number"`
	TopicID            *string `json:"topic_id,omitempty"`
	DeskInstruction    string  `json:"desk_instruction"`
	ExpectedFolderName *string `json:"expected_folder_name,omitempty"`
	ParentID           string  `json:"parent_id,omitempty"`
	Level              int     `json:"level"`
	Sequence           int     `json:"sequence"`
	IsActive           bool    `json:"is_active"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}

// Execute executes the use case
func (uc *CreateWorkPaperItemUseCase) Execute(ctx context.Context, req Request) (*Response, error) {
	// Create service request
	serviceReq := &service.CreateWorkPaperItemRequest{
		Type:            req.Type,
		Number:          req.Number,
		DeskInstruction: req.DeskInstruction,
		Level:           req.Level,
	}

	// Handle TopicID if provided
	if req.TopicID != "" {
		topicUUID, err := parseUUID(req.TopicID)
		if err == nil {
			serviceReq.TopicID = &topicUUID
		}
	}

	// Handle ParentID if provided
	if req.ParentID != "" {
		parentUUID, err := parseUUID(req.ParentID)
		if err == nil {
			serviceReq.ParentID = &parentUUID
		}
	}

	// Call service
	item, err := uc.deskService.CreateWorkPaperItem(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	// Convert to response
	response := &Response{
		ID:                 item.ID.String(),
		Type:               item.Type,
		Number:             item.Number,
		DeskInstruction:    item.DeskInstruction,
		ExpectedFolderName: item.ExpectedFolderName,
		Level:              item.Level,
		Sequence:           item.Sequence,
		IsActive:           item.IsActive,
		CreatedAt:          item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:          item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Handle TopicID if present
	if item.TopicID != nil {
		topicIDStr := item.TopicID.String()
		response.TopicID = &topicIDStr
	}

	// Handle ParentID if present
	if item.ParentID != nil {
		response.ParentID = item.ParentID.String()
	}

	return response, nil
}

// Backward compatibility aliases (deprecated)
type (
	CreateMasterLakipItemUseCase = CreateWorkPaperItemUseCase
	CreateRequest                = Request
	CreateResponse               = Response
)

// NewCreateMasterLakipItemUseCase creates a new use case instance (deprecated)
func NewCreateMasterLakipItemUseCase(deskService service.DeskService) *CreateMasterLakipItemUseCase {
	return NewCreateWorkPaperItemUseCase(deskService)
}
