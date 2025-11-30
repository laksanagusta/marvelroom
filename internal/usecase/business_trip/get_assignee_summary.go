package business_trip

import (
	"context"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
)

type GetAssigneeSummaryUseCase struct {
	businessTripRepo repository.BusinessTripRepository
	assigneeRepo     repository.AssigneeRepository
}

func NewGetAssigneeSummaryUseCase(businessTripRepo repository.BusinessTripRepository, assigneeRepo repository.AssigneeRepository) *GetAssigneeSummaryUseCase {
	return &GetAssigneeSummaryUseCase{
		businessTripRepo: businessTripRepo,
		assigneeRepo:     assigneeRepo,
	}
}

func (uc *GetAssigneeSummaryUseCase) Execute(ctx context.Context, assigneeID string) (*AssigneeSummary, error) {
	// Get assignee with their transactions
	assignee, err := uc.assigneeRepo.GetAssigneeByID(ctx, assigneeID)
	if err != nil {
		return nil, err
	}
	if assignee == nil {
		return nil, entity.ErrAssigneeNotFound
	}

	// Get transactions
	transactions, err := uc.businessTripRepo.GetTransactionsByAssigneeID(ctx, assigneeID)
	if err != nil {
		return nil, err
	}

	// Calculate summary
	costByType := make(map[string]float64)

	for _, tx := range transactions {
		costByType[string(tx.GetType())] += tx.GetSubtotal()
	}

	return &AssigneeSummary{
		AssigneeID:        assigneeID,
		AssigneeName:      assignee.GetName(),
		TotalCost:         assignee.GetTotalCost(),
		TotalTransactions: len(transactions),
		CostByType:        costByType,
	}, nil
}
