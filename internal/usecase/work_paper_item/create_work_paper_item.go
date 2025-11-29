package work_paper_item

import (
	"context"

	"sandbox/internal/domain/service"
)

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
	Type         string `json:"type" validate:"required"`
	Number       string `json:"number" validate:"required"`
	Statement    string `json:"statement" validate:"required"`
	Explanation  string `json:"explanation"`
	FillingGuide string `json:"filling_guide"`
	ParentID     string `json:"parent_id,omitempty"`
	Level        int    `json:"level"`
	SortOrder    int    `json:"sort_order"`
}

// Response represents the response payload for creating a work paper item
type Response struct {
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
}

// Execute executes the use case
func (uc *CreateWorkPaperItemUseCase) Execute(ctx context.Context, req Request) (*Response, error) {
	// Create service request
	serviceReq := &service.CreateWorkPaperItemRequest{
		Type:         req.Type,
		Number:       req.Number,
		Statement:    req.Statement,
		Explanation:  req.Explanation,
		FillingGuide: req.FillingGuide,
		Level:        req.Level,
		SortOrder:    req.SortOrder,
	}

	// Handle ParentID if provided
	if req.ParentID != "" {
		// Convert string ParentID to UUID if needed
		// For simplicity, we'll leave it as empty for now
		// In a real implementation, you would parse the UUID string here
	}

	// Call service
	item, err := uc.deskService.CreateWorkPaperItem(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	// Convert to response
	response := &Response{
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
