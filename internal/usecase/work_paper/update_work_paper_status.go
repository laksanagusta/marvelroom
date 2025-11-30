package work_paper

import (
	"context"
	"errors"

	"sandbox/internal/domain/service"
)

// UpdateWorkPaperStatusUseCase handles updating work paper status
type UpdateWorkPaperStatusUseCase struct {
	deskService service.DeskService
}

// NewUpdateWorkPaperStatusUseCase creates a new use case instance
func NewUpdateWorkPaperStatusUseCase(deskService service.DeskService) *UpdateWorkPaperStatusUseCase {
	return &UpdateWorkPaperStatusUseCase{
		deskService: deskService,
	}
}

// ValidStatusTransitions defines allowed status transitions
var ValidStatusTransitions = map[string][]string{
	"draft":         {"ongoing", "draft"},
	"ongoing":       {"ready_to_sign", "ongoing", "draft"},
	"ready_to_sign": {"completed", "ready_to_sign", "ongoing"},
	"completed":     {"completed", "ready_to_sign"}, // Allow reopening if needed
}

// AllValidStatuses defines all valid status values
var AllValidStatuses = []string{"draft", "ongoing", "ready_to_sign", "completed"}

// UpdateStatusRequest represents the request payload for updating work paper status
type UpdateStatusRequest struct {
	ID     string `json:"id" validate:"required"`
	Status string `json:"status" validate:"required,oneof=draft ongoing ready_to_sign completed"`
}

// UpdateStatusResponse represents the response payload for updating work paper status
type UpdateStatusResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Year           int    `json:"year"`
	Semester       int    `json:"semester"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// ValidateStatusTransition checks if the status transition is valid
func ValidateStatusTransition(currentStatus, newStatus string) error {
	if currentStatus == newStatus {
		return nil // Allow same status (no change)
	}

	allowedTransitions, exists := ValidStatusTransitions[currentStatus]
	if !exists {
		return errors.New("invalid current status")
	}

	for _, allowed := range allowedTransitions {
		if allowed == newStatus {
			return nil
		}
	}

	return errors.New("invalid status transition from " + currentStatus + " to " + newStatus)
}

// Execute executes the use case
func (uc *UpdateWorkPaperStatusUseCase) Execute(ctx context.Context, req UpdateStatusRequest) (*UpdateStatusResponse, error) {
	// Get current work paper to validate status transition
	currentWorkPaper, err := uc.deskService.GetWorkPaper(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// Validate status transition
	if err := ValidateStatusTransition(currentWorkPaper.Status, req.Status); err != nil {
		return nil, err
	}

	// Update work paper status
	err = uc.deskService.UpdateWorkPaperStatus(ctx, req.ID, req.Status)
	if err != nil {
		return nil, err
	}

	// Get the updated work paper to return full response
	updatedWorkPaper, err := uc.deskService.GetWorkPaper(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// Convert to response
	response := &UpdateStatusResponse{
		ID:             updatedWorkPaper.ID.String(),
		OrganizationID: updatedWorkPaper.OrganizationID.String(),
		Year:           updatedWorkPaper.Year,
		Semester:       updatedWorkPaper.Semester,
		Status:         updatedWorkPaper.Status,
		CreatedAt:      updatedWorkPaper.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      updatedWorkPaper.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return response, nil
}

// GetStatusTransitions returns the allowed transitions for a given status
func GetStatusTransitions(currentStatus string) []string {
	if transitions, exists := ValidStatusTransitions[currentStatus]; exists {
		return transitions
	}
	return []string{}
}
