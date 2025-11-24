package business_trip

import (
	"context"
	"fmt"

	"sandbox/internal/domain/repository"
)

type DeleteAssigneeUseCase struct {
	businessTripRepo repository.BusinessTripRepository
}

func NewDeleteAssigneeUseCase(businessTripRepo repository.BusinessTripRepository) *DeleteAssigneeUseCase {
	return &DeleteAssigneeUseCase{
		businessTripRepo: businessTripRepo,
	}
}

func (uc *DeleteAssigneeUseCase) Execute(ctx context.Context, assigneeID string) error {
	assignee, err := uc.businessTripRepo.GetAssigneeByID(ctx, assigneeID)
	if err != nil {
		return fmt.Errorf("failed to get assignee: %w", err)
	}
	if assignee == nil {
		return fmt.Errorf("assignee not found")
	}

	err = uc.businessTripRepo.DeleteAssignee(ctx, assigneeID)
	if err != nil {
		return fmt.Errorf("failed to delete assignee: %w", err)
	}

	return nil
}