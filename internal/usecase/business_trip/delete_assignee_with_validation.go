package business_trip

import (
	"context"
	"fmt"

	"sandbox/internal/domain/repository"
)

type DeleteAssigneeWithValidationUseCase struct {
	businessTripRepo repository.BusinessTripRepository
	assigneeRepo     repository.AssigneeRepository
}

func NewDeleteAssigneeWithValidationUseCase(businessTripRepo repository.BusinessTripRepository, assigneeRepo repository.AssigneeRepository) *DeleteAssigneeWithValidationUseCase {
	return &DeleteAssigneeWithValidationUseCase{
		businessTripRepo: businessTripRepo,
		assigneeRepo:     assigneeRepo,
	}
}

type DeleteAssigneeWithValidationRequest struct {
	BusinessTripID string `params:"businessTripId" json:"businessTripId"`
	AssigneeID     string `params:"assigneeId" json:"assigneeId"`
}

func (r DeleteAssigneeWithValidationRequest) Validate() error {
	if r.BusinessTripID == "" {
		return fmt.Errorf("business trip ID is required")
	}
	if r.AssigneeID == "" {
		return fmt.Errorf("assignee ID is required")
	}
	return nil
}

func (uc *DeleteAssigneeWithValidationUseCase) Execute(ctx context.Context, req DeleteAssigneeWithValidationRequest) error {
	// Validate request
	if err := req.Validate(); err != nil {
		return err
	}

	// Verify business trip exists
	_, err := uc.businessTripRepo.GetByID(ctx, req.BusinessTripID)
	if err != nil {
		return fmt.Errorf("business trip not found")
	}

	// Verify assignee exists and belongs to the business trip
	assignee, err := uc.assigneeRepo.GetAssigneeByID(ctx, req.AssigneeID)
	if err != nil {
		return fmt.Errorf("assignee not found")
	}

	if assignee.BusinessTripID != req.BusinessTripID {
		return fmt.Errorf("assignee does not belong to the specified business trip")
	}

	// Delete assignee (this will also cascade delete transactions)
	err = uc.assigneeRepo.DeleteAssignee(ctx, req.AssigneeID)
	if err != nil {
		return fmt.Errorf("failed to delete assignee: %w", err)
	}

	return nil
}
