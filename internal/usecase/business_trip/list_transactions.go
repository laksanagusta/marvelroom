package business_trip

import (
	"context"
	"fmt"

	"sandbox/internal/domain/repository"
)

type ListTransactionsUseCase struct {
	businessTripRepo repository.BusinessTripRepository
}

func NewListTransactionsUseCase(businessTripRepo repository.BusinessTripRepository) *ListTransactionsUseCase {
	return &ListTransactionsUseCase{
		businessTripRepo: businessTripRepo,
	}
}

type ListTransactionsResponse struct {
	AssigneeID   string                 `json:"assigneeId"`
	Transactions []TransactionResponse  `json:"transactions"`
}

func (uc *ListTransactionsUseCase) Execute(ctx context.Context, assigneeID string) (*ListTransactionsResponse, error) {
	// Verify assignee exists
	_, err := uc.businessTripRepo.GetAssigneeByID(ctx, assigneeID)
	if err != nil {
		return nil, fmt.Errorf("assignee not found")
	}

	// Get transactions for the assignee
	transactions, err := uc.businessTripRepo.GetTransactionsByAssigneeID(ctx, assigneeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	// Convert to response format
	transactionResponses := make([]TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		transactionResponses[i] = TransactionResponse{
			ID:              transaction.ID,
			Name:            transaction.Name,
			Type:            string(transaction.Type),
			Subtype:         string(transaction.Subtype),
			Amount:          transaction.Amount,
			TotalNight:      transaction.TotalNight,
			Subtotal:        transaction.Subtotal,
			Description:     transaction.Description,
			TransportDetail: transaction.TransportDetail,
			CreatedAt:       transaction.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:       transaction.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &ListTransactionsResponse{
		AssigneeID:   assigneeID,
		Transactions: transactionResponses,
	}, nil
}