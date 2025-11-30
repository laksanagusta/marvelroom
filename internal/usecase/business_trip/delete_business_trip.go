package business_trip

import (
	"context"

	"sandbox/internal/domain/repository"
)

type DeleteBusinessTripUseCase struct {
	businessTripRepo repository.BusinessTripRepository
}

func NewDeleteBusinessTripUseCase(businessTripRepo repository.BusinessTripRepository) *DeleteBusinessTripUseCase {
	return &DeleteBusinessTripUseCase{
		businessTripRepo: businessTripRepo,
	}
}

func (uc *DeleteBusinessTripUseCase) Execute(ctx context.Context, id string) error {
	businessTrip, err := uc.businessTripRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if businessTrip == nil {
		return nil
	}

	return uc.businessTripRepo.Delete(ctx, id)
}
