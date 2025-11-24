package business_trip

import (
	"context"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
)

type GetBusinessTripSummaryUseCase struct {
	businessTripRepo repository.BusinessTripRepository
}

func NewGetBusinessTripSummaryUseCase(businessTripRepo repository.BusinessTripRepository) *GetBusinessTripSummaryUseCase {
	return &GetBusinessTripSummaryUseCase{
		businessTripRepo: businessTripRepo,
	}
}

func (uc *GetBusinessTripSummaryUseCase) Execute(ctx context.Context, businessTripID string) (*BusinessTripSummary, error) {
	businessTrip, err := uc.businessTripRepo.GetByID(ctx, businessTripID)
	if err != nil {
		return nil, err
	}
	if businessTrip == nil {
		return nil, entity.ErrBusinessTripNotFound
	}

	assignees, err := uc.businessTripRepo.GetAssigneesByBusinessTripID(ctx, businessTripID)
	if err != nil {
		return nil, err
	}

	costByType := make(map[string]float64)
	totalTransactions := 0

	for _, assignee := range assignees {
		transactions, err := uc.businessTripRepo.GetTransactionsByAssigneeID(ctx, assignee.GetID())
		if err != nil {
			return nil, err
		}

		for _, tx := range transactions {
			costByType[string(tx.GetType())] += tx.GetSubtotal()
			totalTransactions++
		}
	}

	return &BusinessTripSummary{
		BusinessTripID:    businessTripID,
		TotalCost:         businessTrip.GetTotalCost(),
		TotalAssignees:    len(assignees),
		TotalTransactions: totalTransactions,
		CostByType:        costByType,
	}, nil
}
