package work_paper_item

import (
	"context"

	"sandbox/internal/domain/repository"
	"sandbox/pkg/pagination"
)

// ListWorkPaperItemsUseCase handles listing work paper items
type ListWorkPaperItemsUseCase struct {
	workPaperItemRepo repository.WorkPaperItemRepository
}

// NewListWorkPaperItemsUseCase creates a new use case instance
func NewListWorkPaperItemsUseCase(workPaperItemRepo repository.WorkPaperItemRepository) *ListWorkPaperItemsUseCase {
	return &ListWorkPaperItemsUseCase{
		workPaperItemRepo: workPaperItemRepo,
	}
}

// ListRequest represents the request payload for listing work paper items
type ListRequest struct {
	Search   string `json:"search"`
	Type     string `json:"type"`
	IsActive *bool  `json:"is_active"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

// ListResponse represents the response payload for listing work paper items
type ListResponse struct {
	Data     []ItemResponse `json:"data"`
	Metadata Metadata       `json:"metadata"`
}

// Metadata represents pagination metadata
type Metadata struct {
	Count       int `json:"count"`
	TotalCount  int `json:"total_count"`
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	PageSize    int `json:"page_size"`
}

// ItemResponse represents a single work paper item in the response
type ItemResponse struct {
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
func (uc *ListWorkPaperItemsUseCase) Execute(ctx context.Context, params *pagination.QueryParams) ([]*ItemResponse, *pagination.PagedResponse, error) {
	workPaperItems, totalCount, err := uc.workPaperItemRepo.List(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	// Convert entities to response DTOs
	var responses []*ItemResponse
	for _, item := range workPaperItems {
		response := &ItemResponse{
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

		responses = append(responses, response)
	}

	totalPages := int(totalCount) / params.Pagination.Limit
	if int(totalCount)%params.Pagination.Limit > 0 {
		totalPages++
	}

	return responses, &pagination.PagedResponse{
		Page:       params.Pagination.Page,
		Limit:      params.Pagination.Limit,
		TotalItems: totalCount,
		TotalPages: totalPages,
	}, nil
}

// Backward compatibility aliases (deprecated)
type (
	ListMasterLakipItemsUseCase = ListWorkPaperItemsUseCase
	ListRequestLegacy          = ListRequest
	ListResponseLegacy         = ListResponse
	ItemResponseLegacy         = ItemResponse
)

// NewListMasterLakipItemsUseCase creates a new use case instance (deprecated)
func NewListMasterLakipItemsUseCase(workPaperItemRepo repository.WorkPaperItemRepository) *ListMasterLakipItemsUseCase {
	return NewListWorkPaperItemsUseCase(workPaperItemRepo)
}