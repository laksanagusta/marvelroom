package business_trip

import (
	"context"
	"fmt"
	"time"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/pkg/database"
)

// VerifyBusinessTripRequest represents the request to verify a business trip
type VerifyBusinessTripRequest struct {
	BusinessTripID     string `params:"tripId" json:"tripId"`
	VerificationStatus string `json:"status"`             // "approved" or "rejected"
	VerificationNotes  string `json:"verification_notes"` // Optional notes
}

func (r VerifyBusinessTripRequest) Validate() error {
	// Validate BusinessTripID
	if r.BusinessTripID == "" {
		return fmt.Errorf("business trip ID is required")
	}

	// Validate status
	validStatuses := map[string]bool{
		"approved": true,
		"rejected": true,
	}
	if !validStatuses[r.VerificationStatus] {
		return fmt.Errorf("status must be one of: approved, rejected")
	}

	return nil
}

// VerifyBusinessTripResponse represents the response after verification
type VerifyBusinessTripResponse struct {
	ID                 string  `json:"id"`
	BusinessTripID     string  `json:"business_trip_id"`
	UserID             string  `json:"user_id"`
	UserName           string  `json:"user_name"`
	EmployeeNumber     string  `json:"employee_number"`
	Position           string  `json:"position"`
	Status             string  `json:"status"`
	VerifiedAt         *string `json:"verified_at"`
	VerificationNotes  string  `json:"verification_notes"`
	BusinessTripStatus string  `json:"business_trip_status"` // Updated business trip status
}

type VerifyBusinessTripUseCase struct {
	businessTripRepo repository.BusinessTripRepository
	db               database.DB
}

// getUserIDFromContext extracts user ID from context
// This assumes the middleware sets the user ID in context
func getUserIDFromContext(ctx context.Context) (string, error) {
	userID, exists := ctx.Value("user_id").(string)
	if !exists || userID == "" {
		return "", fmt.Errorf("user not authenticated or user_id not found in context")
	}
	return userID, nil
}

func NewVerifyBusinessTripUseCase(businessTripRepo repository.BusinessTripRepository, userService interface{}, db database.DB) *VerifyBusinessTripUseCase {
	return &VerifyBusinessTripUseCase{
		businessTripRepo: businessTripRepo,
		db:               db,
	}
}

func (uc *VerifyBusinessTripUseCase) Execute(ctx context.Context, req VerifyBusinessTripRequest, authenticatedUser entity.AuthenticatedUser) (*VerifyBusinessTripResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	var result *VerifyBusinessTripResponse
	err := uc.db.WithTransaction(ctx, func(ctx context.Context, tx database.DBTx) error {
		// Create transaction-aware repository
		businessTripRepoWithTx := uc.businessTripRepo.(interface {
			WithTransaction(database.DBTx) repository.BusinessTripRepository
		}).WithTransaction(tx)

		// Get verificator for this business trip and user
		verificator, err := businessTripRepoWithTx.GetVerificatorByBusinessTripIDAndUserID(ctx, req.BusinessTripID, authenticatedUser.ID)
		if err != nil {
			return fmt.Errorf("failed to get verificator: %w", err)
		}

		// Check if verificator is still pending
		if !verificator.IsPending() {
			return fmt.Errorf("verificator has already %s this business trip", verificator.GetStatus())
		}

		// Get business trip by ID
		businessTrip, err := businessTripRepoWithTx.GetByID(ctx, req.BusinessTripID)
		if err != nil {
			return fmt.Errorf("failed to get business trip: %w", err)
		}

		// Check if business trip is in ready_to_verify status
		if businessTrip.GetStatus() != entity.BusinessTripStatusReadyToVerify {
			return fmt.Errorf("business trip must be in ready_to_verify status to be verified, current status: %s", businessTrip.GetStatus())
		}

		// Update verificator status
		verificatorStatus := entity.VerificatorStatus(req.VerificationStatus)
		if err := verificator.UpdateStatus(verificatorStatus, req.VerificationNotes); err != nil {
			return fmt.Errorf("failed to update verificator status: %w", err)
		}

		// Update verificator in database
		updatedVerificator, err := businessTripRepoWithTx.UpdateVerificator(ctx, verificator)
		if err != nil {
			return fmt.Errorf("failed to update verificator: %w", err)
		}

		// Check if all verificators have now responded (approved or rejected)
		allVerificators, err := businessTripRepoWithTx.GetVerificatorsByBusinessTripID(ctx, req.BusinessTripID)
		if err != nil {
			return fmt.Errorf("failed to get all verificators: %w", err)
		}

		// Update business trip status based on verificator responses
		newBusinessTripStatus := businessTrip.GetStatus()
		allApproved := true
		anyRejected := false

		for _, v := range allVerificators {
			if v.IsRejected() {
				anyRejected = true
				allApproved = false
				break
			}
			if v.IsPending() {
				allApproved = false
				break
			}
		}

		if anyRejected {
			// If any verificator rejected, mark as ongoing (so it can be fixed and resubmitted)
			newBusinessTripStatus = entity.BusinessTripStatusOngoing
		} else if allApproved {
			// If all approved, mark as ongoing (ready for execution)
			newBusinessTripStatus = entity.BusinessTripStatusOngoing
		}

		// Update business trip status if changed
		if newBusinessTripStatus != businessTrip.GetStatus() {
			if err := businessTrip.UpdateStatus(newBusinessTripStatus); err != nil {
				return fmt.Errorf("failed to update business trip status: %w", err)
			}

			// Update business trip in database
			_, err = businessTripRepoWithTx.Update(ctx, businessTrip)
			if err != nil {
				return fmt.Errorf("failed to update business trip: %w", err)
			}
		}

		// Create response
		var verifiedAt *string
		if updatedVerificator.GetVerifiedAt() != nil {
			verified := updatedVerificator.GetVerifiedAt().Format(time.RFC3339)
			verifiedAt = &verified
		}

		result = &VerifyBusinessTripResponse{
			ID:                 updatedVerificator.GetID(),
			BusinessTripID:     updatedVerificator.GetBusinessTripID(),
			UserID:             authenticatedUser.ID, // Use the authenticated user ID
			UserName:           verificator.UserName, // Use the name from original verificator
			EmployeeNumber:     verificator.EmployeeNumber,
			Position:           verificator.Position,
			Status:             string(updatedVerificator.GetStatus()),
			VerifiedAt:         verifiedAt,
			VerificationNotes:  updatedVerificator.GetVerificationNotes(),
			BusinessTripStatus: string(newBusinessTripStatus),
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
