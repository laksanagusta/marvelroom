package business_trip

import (
	"context"

	"sandbox/internal/domain/repository"
	"sandbox/pkg/pagination"
)

type ListBusinessTripsUseCase struct {
	businessTripRepo repository.BusinessTripRepository
}

func NewListBusinessTripsUseCase(businessTripRepo repository.BusinessTripRepository) *ListBusinessTripsUseCase {
	return &ListBusinessTripsUseCase{
		businessTripRepo: businessTripRepo,
	}
}

func (uc *ListBusinessTripsUseCase) Execute(ctx context.Context, params *pagination.QueryParams) ([]*BusinessTripResponse, *pagination.PagedResponse, error) {
	businessTrips, totalCount, err := uc.businessTripRepo.List(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	// Convert entities to response DTOs
	var responses []*BusinessTripResponse
	for _, bt := range businessTrips {
		responses = append(responses, FromEntity(bt))
	}

	totalPages := int(totalCount) / params.Pagination.Limit
	if int(totalCount)%params.Pagination.Limit > 0 {
		totalPages++
	}

	return responses, &pagination.PagedResponse{
		Page:       params.Pagination.Page,
		Limit:      params.Pagination.Limit,
		TotalItems: totalCount,
		TotalPages: totalPages,
	}, nil
}
