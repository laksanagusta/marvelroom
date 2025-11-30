package business_trip

import (
	"context"
	"fmt"

	"sandbox/internal/domain/repository"
)

type DeleteTransactionUseCase struct {
	businessTripRepo repository.BusinessTripRepository
}

func NewDeleteTransactionUseCase(businessTripRepo repository.BusinessTripRepository) *DeleteTransactionUseCase {
	return &DeleteTransactionUseCase{
		businessTripRepo: businessTripRepo,
	}
}

func (uc *DeleteTransactionUseCase) Execute(ctx context.Context, transactionID string) error {
	transaction, err := uc.businessTripRepo.GetTransactionByID(ctx, transactionID)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}
	if transaction == nil {
		return fmt.Errorf("transaction not found")
	}

	err = uc.businessTripRepo.DeleteTransaction(ctx, transactionID)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	return nil
}
