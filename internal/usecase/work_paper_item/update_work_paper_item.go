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
	ID           string `json:"id" validate:"required"`
	Type         string `json:"type" validate:"required"`
	Number       string `json:"number" validate:"required"`
	Statement    string `json:"statement" validate:"required"`
	Explanation  string `json:"explanation"`
	FillingGuide string `json:"filling_guide"`
	ParentID     string `json:"parent_id,omitempty"`
	Level        int    `json:"level"`
	SortOrder    int    `json:"sort_order"`
	IsActive     *bool  `json:"is_active,omitempty"`
}

// UpdateResponse represents the response payload for updating a work paper item
type UpdateResponse struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Number       string `json:"number"`
	Statement    string `json:"statement"`
	Explanation  string `json:"explanation"`
	FillingGuide string `json:"filling_guide"`
	ParentID     string `json:"parent_id,omitempty"`
	Level        int    `json:"level"`
	SortOrder    int    `json:"sort_order"`
	IsActive     bool   `json:"is_active"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	DeletedAt    string `json:"deleted_at,omitempty"`
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

	// Create service request
	serviceReq := &service.UpdateWorkPaperItemRequest{
		ID:           itemID,
		Type:         req.Type,
		Number:       req.Number,
		Statement:    req.Statement,
		Explanation:  req.Explanation,
		FillingGuide: req.FillingGuide,
		ParentID:     parentID,
		Level:        req.Level,
		SortOrder:    &req.SortOrder,
		IsActive:     req.IsActive,
	}

	// Call service
	item, err := uc.deskService.UpdateWorkPaperItem(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	// Convert to response
	response := &UpdateResponse{
		ID:           item.ID.String(),
		Type:         item.Type,
		Number:       item.Number,
		Statement:    item.Statement,
		Explanation:  item.Explanation,
		FillingGuide: item.FillingGuide,
		Level:        item.Level,
		SortOrder:    item.SortOrder,
		IsActive:     item.IsActive,
		CreatedAt:    item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
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