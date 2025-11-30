package business_trip

import (
	"context"
	"fmt"

	"sandbox/internal/domain/repository"
)

type DeleteAssigneeUseCase struct {
	businessTripRepo repository.BusinessTripRepository
	assigneeRepo     repository.AssigneeRepository
}

func NewDeleteAssigneeUseCase(businessTripRepo repository.BusinessTripRepository, assigneeRepo repository.AssigneeRepository) *DeleteAssigneeUseCase {
	return &DeleteAssigneeUseCase{
		businessTripRepo: businessTripRepo,
		assigneeRepo:     assigneeRepo,
	}
}

func (uc *DeleteAssigneeUseCase) Execute(ctx context.Context, assigneeID string) error {
	assignee, err := uc.assigneeRepo.GetAssigneeByID(ctx, assigneeID)
	if err != nil {
		return fmt.Errorf("failed to get assignee: %w", err)
	}
	if assignee == nil {
		return fmt.Errorf("assignee not found")
	}

	err = uc.assigneeRepo.DeleteAssignee(ctx, assigneeID)
	if err != nil {
		return fmt.Errorf("failed to delete assignee: %w", err)
	}

	return nil
}
