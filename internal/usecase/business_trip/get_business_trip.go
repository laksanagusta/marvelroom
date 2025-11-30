package business_trip

import (
	"context"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
)

type GetBusinessTripUseCase struct {
	businessTripRepo repository.BusinessTripRepository
}

func NewGetBusinessTripUseCase(businessTripRepo repository.BusinessTripRepository) *GetBusinessTripUseCase {
	return &GetBusinessTripUseCase{
		businessTripRepo: businessTripRepo,
	}
}

func (uc *GetBusinessTripUseCase) Execute(ctx context.Context, id string) (*BusinessTripResponse, error) {
	businessTrip, err := uc.businessTripRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if businessTrip == nil {
		return nil, entity.ErrBusinessTripNotFound
	}

	return FromEntity(businessTrip), nil
}
