package work_paper

import (
	"context"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/service"
	"sandbox/pkg/pagination"
)

// ListWorkPapersUseCase handles listing work papers
type ListWorkPapersUseCase struct {
	deskService service.DeskService
}

// NewListWorkPapersUseCase creates a new use case instance
func NewListWorkPapersUseCase(deskService service.DeskService) *ListWorkPapersUseCase {
	return &ListWorkPapersUseCase{
		deskService: deskService,
	}
}

// WorkPaperResponse represents a single work paper in the response
type WorkPaperResponse struct {
	ID             string                `json:"id"`
	OrganizationID string                `json:"organization_id"`
	Organization   *OrganizationResponse `json:"organization,omitempty"`
	Year           int                   `json:"year"`
	Semester       int                   `json:"semester"`
	Status         string                `json:"status"`
	CreatedAt      string                `json:"created_at"`
	UpdatedAt      string                `json:"updated_at"`
}

// Execute executes the use case with pagination params
func (uc *ListWorkPapersUseCase) Execute(ctx context.Context, params *pagination.QueryParams) ([]*WorkPaperResponse, *pagination.PagedResponse, error) {
	// Get work papers from service
	workPapers, totalCount, err := uc.deskService.ListWorkPapers(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	// Convert entities to response DTOs
	responses := make([]*WorkPaperResponse, 0, len(workPapers))
	for _, workPaper := range workPapers {
		response := &WorkPaperResponse{
			ID:             workPaper.ID.String(),
			OrganizationID: workPaper.OrganizationID.String(),
			Year:           workPaper.Year,
			Semester:       workPaper.Semester,
			Status:         workPaper.Status,
			CreatedAt:      workPaper.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:      workPaper.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		// Add organization data if available
		if workPaper.Organization != nil {
			response.Organization = &OrganizationResponse{
				ID:      workPaper.Organization.ID.String(),
				Name:    workPaper.Organization.Name,
				Address: workPaper.Organization.Address,
				Type:    workPaper.Organization.Type,
			}
		}

		responses = append(responses, response)
	}

	// Calculate pagination
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

// ToEntity helper function to convert response to entity (if needed)
func (r *WorkPaperResponse) ToEntity() *entity.WorkPaper {
	return nil // Placeholder, implement if needed
}
