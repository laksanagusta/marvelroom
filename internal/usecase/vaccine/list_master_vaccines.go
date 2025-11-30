package vaccine

import (
	"context"

	"sandbox/pkg/pagination"
)

type ListMasterVaccinesRequest struct {
	QueryParams *pagination.QueryParams `query:"page"`
}

type ListMasterVaccinesResponse struct {
	Data       interface{} `json:"data"`
	Message    string      `json:"message"`
	Success    bool        `json:"success"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
}

type ListMasterVaccinesUseCase struct {
	vaccinesRepo VaccinesRepository
}

func NewListMasterVaccinesUseCase(vaccinesRepo VaccinesRepository) *ListMasterVaccinesUseCase {
	return &ListMasterVaccinesUseCase{
		vaccinesRepo: vaccinesRepo,
	}
}

func (uc *ListMasterVaccinesUseCase) Execute(ctx context.Context, req *ListMasterVaccinesRequest) (*ListMasterVaccinesResponse, error) {
	vaccines, total, err := uc.vaccinesRepo.ListMasterVaccines(ctx, req.QueryParams)
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

	return &ListMasterVaccinesResponse{
		Data:       vaccines,
		Message:    "Master vaccines retrieved successfully",
		Success:    true,
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}, nil
}
