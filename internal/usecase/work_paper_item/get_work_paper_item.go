package work_paper_item

import (
	"context"

	"sandbox/internal/domain/repository"
)

// GetWorkPaperItemUseCase handles getting a single work paper item by ID
type GetWorkPaperItemUseCase struct {
	workPaperItemRepo repository.WorkPaperItemRepository
}

// NewGetWorkPaperItemUseCase creates a new use case instance
func NewGetWorkPaperItemUseCase(workPaperItemRepo repository.WorkPaperItemRepository) *GetWorkPaperItemUseCase {
	return &GetWorkPaperItemUseCase{
		workPaperItemRepo: workPaperItemRepo,
	}
}

// GetRequest represents the request payload for getting a work paper item
type GetRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}

// GetResponse represents the response payload for getting a work paper item
type GetResponse struct {
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
func (uc *GetWorkPaperItemUseCase) Execute(ctx context.Context, req GetRequest) (*GetResponse, error) {
	// Get work paper item from repository
	item, err := uc.workPaperItemRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// Convert entity to response DTO
	response := &GetResponse{
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