package vaccine

import (
	"context"

	"sandbox/pkg/pagination"
)

type ListCountriesRequest struct {
	QueryParams *pagination.QueryParams `query:"page"`
}

type ListCountriesResponse struct {
	Data       interface{} `json:"data"`
	Message    string      `json:"message"`
	Success    bool        `json:"success"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
}

type ListCountriesUseCase struct {
	vaccinesRepo VaccinesRepository
}

func NewListCountriesUseCase(vaccinesRepo VaccinesRepository) *ListCountriesUseCase {
	return &ListCountriesUseCase{
		vaccinesRepo: vaccinesRepo,
	}
}

func (uc *ListCountriesUseCase) Execute(ctx context.Context, req *ListCountriesRequest) (*ListCountriesResponse, error) {
	countries, total, err := uc.vaccinesRepo.List(ctx, req.QueryParams)
	if err != nil {
		return nil, err
	}

	page := 1
	limit := 10
	if req.QueryParams.Pagination.Page > 0 {
		page = req.QueryParams.Pagination.Page
	}
	if req.QueryParams.Pagination.Limit > 0 {
		limit = req.QueryParams.Pagination.Limit
	}

	return &ListCountriesResponse{
		Data:       countries,
		Message:    "Countries retrieved successfully",
		Success:    true,
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}, nil
}
