package business_trip

import (
	"context"

	"sandbox/internal/domain/repository"
)

// GetHistoriesUseCase handles retrieving history for business trips
type GetHistoriesUseCase struct {
	historyRepo repository.BusinessTripHistoryRepository
}

// NewGetHistoriesUseCase creates a new instance of GetHistoriesUseCase
func NewGetHistoriesUseCase(historyRepo repository.BusinessTripHistoryRepository) *GetHistoriesUseCase {
	return &GetHistoriesUseCase{
		historyRepo: historyRepo,
	}
}

// Execute retrieves all history records for a business trip
func (uc *GetHistoriesUseCase) Execute(ctx context.Context, businessTripID string) ([]*BusinessTripHistoryResponse, error) {
	histories, err := uc.historyRepo.FindByBusinessTripID(ctx, businessTripID)
	if err != nil {
		return nil, err
	}

	// Convert to response
	responses := make([]*BusinessTripHistoryResponse, 0, len(histories))
	for _, history := range histories {
		responses = append(responses, &BusinessTripHistoryResponse{
			ID:             history.GetID(),
			BusinessTripID: history.GetBusinessTripID(),
			ChangeType:     string(history.GetChangeType()),
			FieldName:      history.GetFieldName(),
			OldValue:       history.GetOldValue(),
			NewValue:       history.GetNewValue(),
			ChangedBy:      history.GetChangedBy(),
			Notes:          history.GetNotes(),
			CreatedAt:      history.GetCreatedAt(),
		})
	}

	return responses, nil
}
