package work_paper_item

import (
	"context"

	"sandbox/internal/domain/service"

	"github.com/google/uuid"
)

// DeleteWorkPaperItemUseCase handles the deletion of work paper items
type DeleteWorkPaperItemUseCase struct {
	deskService service.DeskService
}

// NewDeleteWorkPaperItemUseCase creates a new use case instance
func NewDeleteWorkPaperItemUseCase(deskService service.DeskService) *DeleteWorkPaperItemUseCase {
	return &DeleteWorkPaperItemUseCase{
		deskService: deskService,
	}
}

// DeleteRequest represents the request payload for deleting a work paper item
type DeleteRequest struct {
	ID string `json:"id" validate:"required"`
}

// DeleteResponse represents the response payload for deleting a work paper item
type DeleteResponse struct {
	ID        string `json:"id"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	DeletedAt string `json:"deleted_at"`
}

// Execute executes the use case for deleting a work paper item
func (uc *DeleteWorkPaperItemUseCase) Execute(ctx context.Context, req DeleteRequest) (*DeleteResponse, error) {
	// Parse ID
	itemID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, err
	}

	// Call service to delete the work paper item
	item, err := uc.deskService.DeleteWorkPaperItem(ctx, itemID)
	if err != nil {
		return nil, err
	}

	// Convert to response
	response := &DeleteResponse{
		ID:        item.ID.String(),
		Success:   true,
		Message:   "Work paper item deleted successfully",
		DeletedAt: item.DeletedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return response, nil
}