package work_paper

import (
	"context"

	"sandbox/internal/domain/service"
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

// ListRequest represents the request payload for listing work papers
type ListRequest struct {
	OrganizationID string `json:"organization_id"`
	Year           *int   `json:"year"`
	Semester       *int   `json:"semester"`
	Status         string `json:"status"`
	Page           int    `json:"page" validate:"min=1"`
	PageSize       int    `json:"page_size" validate:"min=1,max=100"`
}

// ListResponse represents the response payload for listing work papers
type ListResponse struct {
	Data     []WorkPaperResponse `json:"data"`
	Metadata Metadata            `json:"metadata"`
}

// Metadata represents pagination metadata
type Metadata struct {
	Count       int `json:"count"`
	TotalCount  int `json:"total_count"`
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	PageSize    int `json:"page_size"`
}

// WorkPaperResponse represents a single work paper in the response
type WorkPaperResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Year           int    `json:"year"`
	Semester       int    `json:"semester"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// Execute executes the use case
func (uc *ListWorkPapersUseCase) Execute(ctx context.Context, req ListRequest) (*ListResponse, error) {
	// Create service request
	serviceReq := &service.ListWorkPapersRequest{
		OrganizationID: req.OrganizationID,
		Year:           req.Year,
		Semester:       req.Semester,
		Status:         req.Status,
	}

	// Get work papers from service
	workPapers, totalCount, err := uc.deskService.ListWorkPapers(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	// Convert entities to response DTOs
	var responses []WorkPaperResponse
	for _, workPaper := range workPapers {
		response := WorkPaperResponse{
			ID:             workPaper.ID.String(),
			OrganizationID: workPaper.OrganizationID.String(),
			Year:           workPaper.Year,
			Semester:       workPaper.Semester,
			Status:         workPaper.Status,
			CreatedAt:      workPaper.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:      workPaper.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		responses = append(responses, response)
	}

	// Calculate pagination
	page := req.Page
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 10 // default page size
	}

	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize > 0 {
		totalPages++
	}

	return &ListResponse{
		Data: responses,
		Metadata: Metadata{
			Count:       len(responses),
			TotalCount:  int(totalCount),
			CurrentPage: page,
			TotalPage:   totalPages,
			PageSize:    pageSize,
		},
	}, nil
}
