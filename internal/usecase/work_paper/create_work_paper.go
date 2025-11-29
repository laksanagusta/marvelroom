package work_paper

import (
	"context"

	"sandbox/internal/domain/service"
)

// CreateWorkPaperUseCase handles the creation of work paper
type CreateWorkPaperUseCase struct {
	deskService service.DeskService
}

// InjectToPublicUseCase injects this use case into public use cases (for backwards compatibility)
func (uc *CreateWorkPaperUseCase) InjectToPublicUseCase() *CreateWorkPaperUseCase {
	return uc
}

// NewCreateWorkPaperUseCase creates a new use case instance
func NewCreateWorkPaperUseCase(deskService service.DeskService) *CreateWorkPaperUseCase {
	return &CreateWorkPaperUseCase{
		deskService: deskService,
	}
}

// CreateRequest represents the request payload for creating a work paper
type CreateRequest struct {
	OrganizationID string `json:"organization_id" validate:"required"`
	Year           int    `json:"year" validate:"required,min=2000,max=2100"`
	Semester       int    `json:"semester" validate:"required,oneof=1 2"`
}

// CreateResponse represents the response payload for creating a work paper
type CreateResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Year           int    `json:"year"`
	Semester       int    `json:"semester"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// Execute executes the use case
func (uc *CreateWorkPaperUseCase) Execute(ctx context.Context, req CreateRequest) (*CreateResponse, error) {
	// Create service request
	serviceReq := &service.CreateWorkPaperRequest{
		OrganizationID: req.OrganizationID,
		Year:           req.Year,
		Semester:       req.Semester,
	}

	// Call service
	workPaper, err := uc.deskService.CreateWorkPaper(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	// Convert to response
	response := &CreateResponse{
		ID:             workPaper.ID.String(),
		OrganizationID: workPaper.OrganizationID.String(),
		Year:           workPaper.Year,
		Semester:       workPaper.Semester,
		Status:         workPaper.Status,
		CreatedAt:      workPaper.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      workPaper.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return response, nil
}

// Backward compatibility aliases (deprecated)
type (
	CreatePaperWorkUseCase = CreateWorkPaperUseCase
	CreateRequestLegacy  = CreateRequest
	CreateResponseLegacy = CreateResponse
)

// NewCreatePaperWorkUseCase creates a new use case instance (deprecated)
func NewCreatePaperWorkUseCase(deskService service.DeskService) *CreatePaperWorkUseCase {
	return NewCreateWorkPaperUseCase(deskService)
}