package work_paper_item

import (
	"context"

	"sandbox/internal/domain/service"

	"github.com/google/uuid"
)

// UpdateWorkPaperItemUseCase handles the updating of work paper items
type UpdateWorkPaperItemUseCase struct {
	deskService service.DeskService
}

// NewUpdateWorkPaperItemUseCase creates a new use case instance
func NewUpdateWorkPaperItemUseCase(deskService service.DeskService) *UpdateWorkPaperItemUseCase {
	return &UpdateWorkPaperItemUseCase{
		deskService: deskService,
	}
}

// UpdateRequest represents the request payload for updating a work paper item
type UpdateRequest struct {
	ID                 string  `json:"id" validate:"required"`
	Type               string  `json:"type" validate:"required"`
	Number             string  `json:"number" validate:"required"`
	TopicID            string  `json:"topic_id"`
	DeskInstruction    string  `json:"desk_instruction" validate:"required"`
	ExpectedFolderName *string `json:"expected_folder_name,omitempty"`
	ParentID           string  `json:"parent_id,omitempty"`
	Level              int     `json:"level"`
	Sequence           int     `json:"sequence"`
	IsActive           *bool   `json:"is_active,omitempty"`
}

// UpdateResponse represents the response payload for updating a work paper item
type UpdateResponse struct {
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
	DeletedAt          string  `json:"deleted_at,omitempty"`
}

// Execute executes the use case for updating a work paper item
func (uc *UpdateWorkPaperItemUseCase) Execute(ctx context.Context, req UpdateRequest) (*UpdateResponse, error) {
	// Parse ID
	itemID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, err
	}

	// Parse ParentID if provided
	var parentID *uuid.UUID
	if req.ParentID != "" {
		parsedParentID, err := uuid.Parse(req.ParentID)
		if err != nil {
			return nil, err
		}
		parentID = &parsedParentID
	}

	// Parse TopicID if provided
	var topicID *uuid.UUID
	if req.TopicID != "" {
		parsedTopicID, err := uuid.Parse(req.TopicID)
		if err != nil {
			return nil, err
		}
		topicID = &parsedTopicID
	}

	// Create service request
	sequence := req.Sequence
	serviceReq := &service.UpdateWorkPaperItemRequest{
		ID:                 itemID,
		Type:               req.Type,
		Number:             req.Number,
		TopicID:            topicID,
		DeskInstruction:    req.DeskInstruction,
		ExpectedFolderName: req.ExpectedFolderName,
		ParentID:           parentID,
		Level:              req.Level,
		Sequence:           &sequence,
		IsActive:           req.IsActive,
	}

	// Call service
	item, err := uc.deskService.UpdateWorkPaperItem(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	// Convert to response
	response := &UpdateResponse{
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

	// Handle DeletedAt if present
	if item.DeletedAt != nil {
		response.DeletedAt = item.DeletedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	return response, nil
}
