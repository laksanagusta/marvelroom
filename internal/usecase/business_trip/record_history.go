package business_trip

import (
	"context"
	"fmt"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
)

// RecordHistoryUseCase handles recording history for business trips
type RecordHistoryUseCase struct {
	historyRepo repository.BusinessTripHistoryRepository
}

// NewRecordHistoryUseCase creates a new instance of RecordHistoryUseCase
func NewRecordHistoryUseCase(historyRepo repository.BusinessTripHistoryRepository) *RecordHistoryUseCase {
	return &RecordHistoryUseCase{
		historyRepo: historyRepo,
	}
}

// RecordHistoryInput represents the input for recording history
type RecordHistoryInput struct {
	BusinessTripID string
	ChangeType     entity.BusinessTripHistoryChangeType
	FieldName      string
	OldValue       string
	NewValue       string
	ChangedBy      string
	Notes          string
}

// Execute records a history entry for a business trip
func (uc *RecordHistoryUseCase) Execute(ctx context.Context, input RecordHistoryInput) error {
	// Create new history record
	history, err := entity.NewBusinessTripHistory(input.BusinessTripID, input.ChangeType)
	if err != nil {
		return fmt.Errorf("failed to create history entity: %w", err)
	}

	// Set optional fields
	if input.FieldName != "" {
		history.SetFieldName(input.FieldName)
	}

	if input.OldValue != "" {
		history.SetOldValue(input.OldValue)
	}

	if input.NewValue != "" {
		history.SetNewValue(input.NewValue)
	}

	if input.ChangedBy != "" {
		history.SetChangedBy(input.ChangedBy)
	}

	if input.Notes != "" {
		history.SetNotes(input.Notes)
	}

	// Save to repository
	if err := uc.historyRepo.Create(ctx, history); err != nil {
		return fmt.Errorf("failed to save history: %w", err)
	}

	return nil
}
