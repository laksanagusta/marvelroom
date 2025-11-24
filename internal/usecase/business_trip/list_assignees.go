package business_trip

import (
	"context"
	"fmt"

	"sandbox/internal/domain/repository"
)

type ListAssigneesUseCase struct {
	businessTripRepo repository.BusinessTripRepository
}

func NewListAssigneesUseCase(businessTripRepo repository.BusinessTripRepository) *ListAssigneesUseCase {
	return &ListAssigneesUseCase{
		businessTripRepo: businessTripRepo,
	}
}

type ListAssigneesResponse struct {
	BusinessTripID string                 `json:"businessTripId"`
	Assignees      []AssigneeResponse    `json:"assignees"`
}

func (uc *ListAssigneesUseCase) Execute(ctx context.Context, businessTripID string) (*ListAssigneesResponse, error) {
	// Verify business trip exists
	_, err := uc.businessTripRepo.GetByID(ctx, businessTripID)
	if err != nil {
		return nil, fmt.Errorf("business trip not found")
	}

	// Get assignees for the business trip
	assignees, err := uc.businessTripRepo.GetAssigneesByBusinessTripID(ctx, businessTripID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignees: %w", err)
	}

	// Convert to response format
	assigneeResponses := make([]AssigneeResponse, len(assignees))
	for i, assignee := range assignees {
		assigneeResponses[i] = AssigneeResponse{
			ID:         assignee.ID,
			Name:       assignee.Name,
			SPDNumber:  assignee.SPDNumber,
			EmployeeID: assignee.EmployeeID,
			Position:   assignee.Position,
			Rank:       assignee.Rank,
			TotalCost:  assignee.GetTotalCost(),
			CreatedAt:  assignee.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  assignee.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		// Convert transactions
		transactionResponses := make([]TransactionResponse, len(assignee.Transactions))
		for j, transaction := range assignee.Transactions {
			transactionResponses[j] = TransactionResponse{
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
		assigneeResponses[i].Transactions = transactionResponses
	}

	return &ListAssigneesResponse{
		BusinessTripID: businessTripID,
		Assignees:      assigneeResponses,
	}, nil
}