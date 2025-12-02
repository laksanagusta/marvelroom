package work_paper

import (
	"context"

	"sandbox/internal/domain/service"
)

// GetWorkPaperDetailsUseCase handles getting work paper with all related details
type GetWorkPaperDetailsUseCase struct {
	deskService service.DeskService
}

// NewGetWorkPaperDetailsUseCase creates a new use case instance
func NewGetWorkPaperDetailsUseCase(deskService service.DeskService) *GetWorkPaperDetailsUseCase {
	return &GetWorkPaperDetailsUseCase{
		deskService: deskService,
	}
}

// GetWorkPaperDetailsResponse represents the detailed response with all related data
type GetWorkPaperDetailsResponse struct {
	ID             string                `json:"id"`
	OrganizationID string                `json:"organization_id"`
	Organization   *OrganizationResponse `json:"organization,omitempty"`
	Year           int                   `json:"year"`
	Semester       int                   `json:"semester"`
	Status         string                `json:"status"`
	CreatedAt      string                `json:"created_at"`
	UpdatedAt      string                `json:"updated_at"`
	// Include related data
	WorkPaperNotes []*WorkPaperNoteResponse      `json:"work_paper_notes,omitempty"`
	Signatures     []*WorkPaperSignatureResponse `json:"signatures,omitempty"`
}

// WorkPaperNoteResponse represents a work paper note in the detailed response
type WorkPaperNoteResponse struct {
	ID           string `json:"id"`
	WorkPaperID  string `json:"work_paper_id"`
	Statement    string `json:"statement"`
	Explanation  string `json:"explanation"`
	FillingGuide string `json:"filling_guide"`
	Status       string `json:"status"`
	DriveLink    string `json:"gdrive_link"`
	IsValid      *bool  `json:"is_valid"`
	Notes        string `json:"notes"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// WorkPaperSignatureResponse represents a work paper signature in the detailed response
type WorkPaperSignatureResponse struct {
	ID            string `json:"id"`
	WorkPaperID   string `json:"work_paper_id"`
	UserID        string `json:"user_id"`
	UserName      string `json:"user_name"`
	UserEmail     string `json:"user_email,omitempty"`
	UserRole      string `json:"user_role"`
	SignatureType string `json:"signature_type"`
	Status        string `json:"status"`
	Notes         string `json:"notes,omitempty"`
	SignedAt      string `json:"signed_at,omitempty"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// Execute executes the use case
func (uc *GetWorkPaperDetailsUseCase) Execute(ctx context.Context, workPaperID string) (*GetWorkPaperDetailsResponse, error) {
	// Get the work paper
	workPaper, err := uc.deskService.GetWorkPaper(ctx, workPaperID)
	if err != nil {
		return nil, err
	}

	// Get work paper notes
	notes, err := uc.deskService.GetWorkPaperNotes(ctx, workPaperID)
	if err != nil {
		return nil, err
	}

	// Get work paper signatures
	signatures, err := uc.deskService.GetWorkPaperSignatures(ctx, workPaperID)
	if err != nil {
		return nil, err
	}

	// Convert work paper notes to response format
	var noteResponses []*WorkPaperNoteResponse
	for _, note := range notes {
		noteResponse := &WorkPaperNoteResponse{
			ID:          note.ID.String(),
			WorkPaperID: note.WorkPaperID.String(),
			DriveLink:   note.GetGDriveLink(),
			CreatedAt:   note.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   note.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		if note.IsValid != nil {
			noteResponse.IsValid = note.IsValid
		}

		// Add statement, explanation, and filling guide from master item if available
		if note.MasterItem != nil {
			noteResponse.Statement = note.MasterItem.Statement
			noteResponse.Explanation = note.MasterItem.Explanation
			noteResponse.FillingGuide = note.MasterItem.FillingGuide
			noteResponse.Status = "active"
			if note.MasterItem.DeletedAt != nil {
				noteResponse.Status = "inactive"
			}
		}

		noteResponse.Notes = note.GetNotes()

		noteResponses = append(noteResponses, noteResponse)
	}

	// Convert work paper signatures to response format
	var signatureResponses []*WorkPaperSignatureResponse
	for _, signature := range signatures {
		signatureResponse := &WorkPaperSignatureResponse{
			ID:            signature.ID.String(),
			WorkPaperID:   signature.WorkPaperID.String(),
			UserID:        signature.UserID,
			UserName:      signature.UserName,
			UserEmail:     signature.GetUserEmail(),
			UserRole:      signature.GetUserRole(),
			SignatureType: signature.SignatureType,
			Status:        signature.Status,
			Notes:         signature.GetNotes(),
			SignedAt:      signature.GetSignedAt().Format("2006-01-02T15:04:05Z07:00"),
			CreatedAt:     signature.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     signature.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		signatureResponses = append(signatureResponses, signatureResponse)
	}

	// Build the complete response
	response := &GetWorkPaperDetailsResponse{
		ID:             workPaper.ID.String(),
		OrganizationID: workPaper.OrganizationID.String(),
		Year:           workPaper.Year,
		Semester:       workPaper.Semester,
		Status:         workPaper.Status,
		CreatedAt:      workPaper.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      workPaper.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		WorkPaperNotes: noteResponses,
		Signatures:     signatureResponses,
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

	return response, nil
}
