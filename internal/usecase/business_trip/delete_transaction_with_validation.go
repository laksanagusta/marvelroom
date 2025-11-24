package business_trip

import (
	"context"
	"fmt"

	"sandbox/internal/domain/repository"
)

type DeleteTransactionWithValidationUseCase struct {
	businessTripRepo repository.BusinessTripRepository
}

func NewDeleteTransactionWithValidationUseCase(businessTripRepo repository.BusinessTripRepository) *DeleteTransactionWithValidationUseCase {
	return &DeleteTransactionWithValidationUseCase{
		businessTripRepo: businessTripRepo,
	}
}

type DeleteTransactionWithValidationRequest struct {
	BusinessTripID string `params:"businessTripId" json:"businessTripId"`
	AssigneeID     string `params:"assigneeId" json:"assigneeId"`
	TransactionID  string `params:"transactionId" json:"transactionId"`
}

func (r DeleteTransactionWithValidationRequest) Validate() error {
	if r.BusinessTripID == "" {
		return fmt.Errorf("business trip ID is required")
	}
	if r.AssigneeID == "" {
		return fmt.Errorf("assignee ID is required")
	}
	if r.TransactionID == "" {
		return fmt.Errorf("transaction ID is required")
	}
	return nil
}

func (uc *DeleteTransactionWithValidationUseCase) Execute(ctx context.Context, req DeleteTransactionWithValidationRequest) error {
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
	assignee, err := uc.businessTripRepo.GetAssigneeByID(ctx, req.AssigneeID)
	if err != nil {
		return fmt.Errorf("assignee not found")
	}

	if assignee.BusinessTripID != req.BusinessTripID {
		return fmt.Errorf("assignee does not belong to the specified business trip")
	}

	// Check if transaction exists and belongs to the specified assignee
	transaction, err := uc.businessTripRepo.GetTransactionByID(ctx, req.TransactionID)
	if err != nil {
		return fmt.Errorf("transaction not found")
	}

	if transaction.AssigneeID != req.AssigneeID {
		return fmt.Errorf("transaction does not belong to the specified assignee")
	}

	// Delete transaction
	err = uc.businessTripRepo.DeleteTransaction(ctx, req.TransactionID)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	return nil
}