package business_trip

import (
	"context"
	"fmt"
	"time"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
)

type UpdateBusinessTripUseCase struct {
	businessTripRepo repository.BusinessTripRepository
}

func NewUpdateBusinessTripUseCase(businessTripRepo repository.BusinessTripRepository) *UpdateBusinessTripUseCase {
	return &UpdateBusinessTripUseCase{
		businessTripRepo: businessTripRepo,
	}
}

func (uc *UpdateBusinessTripUseCase) Execute(ctx context.Context, req UpdateBusinessTripRequest) (*BusinessTripResponse, error) {
	businessTrip, err := uc.businessTripRepo.GetByID(ctx, req.BusinessTripID)
	if err != nil {
		return nil, fmt.Errorf("failed to get business trip: %w", err)
	}
	if businessTrip == nil {
		return nil, entity.ErrBusinessTripNotFound
	}

	// Update fields if provided
	if req.StartDate.IsSet() {
		startDate, err := time.Parse("2006-01-02", req.StartDate.String)
		if err != nil {
			return nil, err
		}
		businessTrip.StartDate = startDate
	}
	if req.EndDate.IsSet() {
		endDate, err := time.Parse("2006-01-02", req.EndDate.String)
		if err != nil {
			return nil, err
		}
		businessTrip.EndDate = endDate
	}
	if req.ActivityPurpose.IsSet() {
		businessTrip.ActivityPurpose = req.ActivityPurpose.String
	}
	if req.DestinationCity.IsSet() {
		businessTrip.DestinationCity = req.DestinationCity.String
	}
	if req.SPDDate.IsSet() {
		spdDate, err := time.Parse("2006-01-02", req.SPDDate.String)
		if err != nil {
			return nil, err
		}
		businessTrip.SPDDate = spdDate
	}
	if req.DepartureDate.IsSet() {
		departureDate, err := time.Parse("2006-01-02", req.DepartureDate.String)
		if err != nil {
			return nil, err
		}
		businessTrip.DepartureDate = departureDate
	}
	if req.ReturnDate.IsSet() {
		returnDate, err := time.Parse("2006-01-02", req.ReturnDate.String)
		if err != nil {
			return nil, err
		}
		businessTrip.ReturnDate = returnDate
	}

	// Validate the updated business trip
	if businessTrip.StartDate.After(businessTrip.EndDate) {
		return nil, entity.ErrInvalidDateRange
	}
	if businessTrip.DepartureDate.After(businessTrip.ReturnDate) {
		return nil, entity.ErrInvalidDateRange
	}
	if businessTrip.SPDDate.After(businessTrip.DepartureDate) {
		return nil, entity.ErrInvalidDateRange
	}

	// Update status if provided
	if req.Status.IsSet() {
		newStatus := entity.BusinessTripStatus(req.Status.String)
		if err := businessTrip.UpdateStatus(newStatus); err != nil {
			return nil, err
		}
	}

	// Update document link if provided
	if req.DocumentLink.IsSet() {
		businessTrip.UpdateDocumentLink(req.DocumentLink.String)
	}

	// Save to repository
	updatedBusinessTrip, err := uc.businessTripRepo.Update(ctx, businessTrip)
	if err != nil {
		return nil, err
	}

	return FromEntity(updatedBusinessTrip), nil
}
