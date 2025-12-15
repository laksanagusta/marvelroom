package business_trip

import (
	"context"
	"errors"
	"fmt"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"

	"github.com/invopop/validation"
)

// ErrStatusTransitionNotAllowed is returned when a status transition is not allowed
var ErrStatusTransitionNotAllowed = errors.New("status transition not allowed")

// ErrFinalStatusCannotBeChanged is returned when trying to change a final status (canceled or completed)
var ErrFinalStatusCannotBeChanged = errors.New("status cannot be changed once it is canceled or completed")

// UpdateBusinessTripStatusRequest represents the request body for updating business trip status
type UpdateBusinessTripStatusRequest struct {
	BusinessTripID string `json:"-"`
	Status         string `json:"status"`
}

func (r UpdateBusinessTripStatusRequest) Validate() error {
	// Basic validation
	err := validation.ValidateStruct(&r,
		validation.Field(&r.Status, validation.Required, validation.Length(1, 50)),
	)
	if err != nil {
		return err
	}

	// Validate status value
	validStatuses := map[string]bool{
		"draft":           true,
		"ongoing":         true,
		"ready_to_verify": true,
		"completed":       true,
		"canceled":        true,
	}
	if !validStatuses[r.Status] {
		return validation.NewError("status", "must be one of: draft, ongoing, ready_to_verify, completed, canceled")
	}

	return nil
}

// UpdateBusinessTripStatusUseCase handles updating business trip status
type UpdateBusinessTripStatusUseCase struct {
	businessTripRepo repository.BusinessTripRepository
	historyUseCase   *RecordHistoryUseCase
}

// NewUpdateBusinessTripStatusUseCase creates a new instance of UpdateBusinessTripStatusUseCase
func NewUpdateBusinessTripStatusUseCase(
	businessTripRepo repository.BusinessTripRepository,
	historyUseCase *RecordHistoryUseCase,
) *UpdateBusinessTripStatusUseCase {
	return &UpdateBusinessTripStatusUseCase{
		businessTripRepo: businessTripRepo,
		historyUseCase:   historyUseCase,
	}
}

// ValidStatusTransitions defines the allowed status transitions
// Flow: draft -> canceled/ongoing -> ready_to_verify/canceled -> completed/canceled
// Once canceled or completed, status cannot be changed
var ValidStatusTransitions = map[entity.BusinessTripStatus][]entity.BusinessTripStatus{
	entity.BusinessTripStatusDraft: {
		entity.BusinessTripStatusCanceled,
		entity.BusinessTripStatusOngoing,
	},
	entity.BusinessTripStatusOngoing: {
		entity.BusinessTripStatusCanceled,
		entity.BusinessTripStatusReadyToVerify,
	},
	entity.BusinessTripStatusReadyToVerify: {
		entity.BusinessTripStatusCompleted,
		entity.BusinessTripStatusCanceled,
	},
	// canceled and completed are final states - no transitions allowed
	entity.BusinessTripStatusCanceled:  {},
	entity.BusinessTripStatusCompleted: {},
}

// canTransitionStatus checks if transition from current status to new status is allowed
func canTransitionStatus(currentStatus, newStatus entity.BusinessTripStatus) bool {
	allowedTransitions, exists := ValidStatusTransitions[currentStatus]
	if !exists {
		return false
	}

	for _, allowedStatus := range allowedTransitions {
		if allowedStatus == newStatus {
			return true
		}
	}
	return false
}

// Execute updates the business trip status
func (uc *UpdateBusinessTripStatusUseCase) Execute(
	ctx context.Context,
	req UpdateBusinessTripStatusRequest,
	authenticatedUser entity.AuthenticatedUser,
) (*BusinessTripResponse, error) {
	// Get existing business trip
	businessTrip, err := uc.businessTripRepo.GetByID(ctx, req.BusinessTripID)
	if err != nil {
		return nil, fmt.Errorf("failed to get business trip: %w", err)
	}
	if businessTrip == nil {
		return nil, entity.ErrBusinessTripNotFound
	}

	currentStatus := businessTrip.Status
	newStatus := entity.BusinessTripStatus(req.Status)

	// Check if current status is a final status
	if currentStatus == entity.BusinessTripStatusCanceled || currentStatus == entity.BusinessTripStatusCompleted {
		return nil, fmt.Errorf("%w: current status is '%s'", ErrFinalStatusCannotBeChanged, currentStatus)
	}

	// Validate status transition
	if !canTransitionStatus(currentStatus, newStatus) {
		allowedStatuses := ValidStatusTransitions[currentStatus]
		allowedStatusStrings := make([]string, len(allowedStatuses))
		for i, s := range allowedStatuses {
			allowedStatusStrings[i] = string(s)
		}
		return nil, fmt.Errorf("%w: from '%s' you can only transition to: %v", ErrStatusTransitionNotAllowed, currentStatus, allowedStatusStrings)
	}

	// Additional validation: if transitioning to completed
	if newStatus == entity.BusinessTripStatusCompleted {
		// Check document link is required
		if !businessTrip.DocumentLink.Valid || businessTrip.DocumentLink.String == "" {
			return nil, fmt.Errorf("document link is required when marking business trip as completed")
		}

		// Check all verificators must be approved
		if len(businessTrip.Verificators) > 0 {
			pendingCount := businessTrip.GetPendingVerificatorsCount()
			rejectedCount := businessTrip.GetRejectedVerificatorsCount()

			if rejectedCount > 0 {
				return nil, fmt.Errorf("cannot complete business trip: %d verificator(s) have rejected", rejectedCount)
			}

			if pendingCount > 0 {
				return nil, fmt.Errorf("cannot complete business trip: %d verificator(s) have not yet approved", pendingCount)
			}
		}
	}

	// Record old status for history
	oldStatus := string(businessTrip.Status)

	// Update status
	businessTrip.Status = newStatus

	// Save to repository
	updatedBusinessTrip, err := uc.businessTripRepo.Update(ctx, businessTrip)
	if err != nil {
		return nil, fmt.Errorf("failed to update business trip: %w", err)
	}

	// Record history
	if uc.historyUseCase != nil {
		err = uc.historyUseCase.Execute(ctx, RecordHistoryInput{
			BusinessTripID: businessTrip.ID,
			ChangeType:     entity.HistoryChangeTypeStatusChange,
			FieldName:      "status",
			OldValue:       oldStatus,
			NewValue:       string(newStatus),
			ChangedBy:      authenticatedUser.GetFullName(),
		})
		if err != nil {
			// Log error but don't fail the operation
			fmt.Printf("failed to record history: %v\n", err)
		}
	}

	return FromEntity(updatedBusinessTrip), nil
}
