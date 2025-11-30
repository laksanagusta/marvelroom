package business_trip

import (
	"context"
	"fmt"

	"sandbox/internal/domain/repository"
)

type GetAssigneeUseCase struct {
	assigneeRepo repository.AssigneeRepository
}

func NewGetAssigneeUseCase(assigneeRepo repository.AssigneeRepository) *GetAssigneeUseCase {
	return &GetAssigneeUseCase{
		assigneeRepo: assigneeRepo,
	}
}

type GetAssigneeResponse struct {
	ID             string                `json:"id"`
	BusinessTripID string                `json:"businessTripId"`
	Name           string                `json:"name"`
	SPDNumber      string                `json:"spdNumber"`
	EmployeeID     string                `json:"employeeId"`
	EmployeeName   string                `json:"employeeName"`
	Position       string                `json:"position"`
	Rank           string                `json:"rank"`
	TotalCost      float64               `json:"totalCost"`
	Transactions   []TransactionResponse `json:"transactions"`
	CreatedAt      string                `json:"createdAt"`
	UpdatedAt      string                `json:"updatedAt"`
}

func (uc *GetAssigneeUseCase) Execute(ctx context.Context, assigneeID string) (*GetAssigneeResponse, error) {
	assignee, err := uc.assigneeRepo.GetAssigneeByID(ctx, assigneeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignee: %w", err)
	}
	if assignee == nil {
		return nil, fmt.Errorf("assignee not found")
	}

	transactionResponses := make([]TransactionResponse, len(assignee.Transactions))
	for i, transaction := range assignee.Transactions {
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

	return &GetAssigneeResponse{
		ID:             assignee.ID,
		BusinessTripID: assignee.BusinessTripID,
		Name:           assignee.Name,
		SPDNumber:      assignee.SPDNumber,
		EmployeeID:     assignee.EmployeeID,
		EmployeeName:   assignee.EmployeeName,
		Position:       assignee.Position,
		Rank:           assignee.Rank,
		TotalCost:      assignee.GetTotalCost(),
		Transactions:   transactionResponses,
		CreatedAt:      assignee.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      assignee.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
