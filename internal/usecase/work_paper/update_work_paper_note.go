package work_paper

import (
	"context"

	"sandbox/internal/domain/service"
)

// UpdateWorkPaperNoteUseCase handles updating work paper notes
type UpdateWorkPaperNoteUseCase struct {
	deskService service.DeskService
}

// NewUpdateWorkPaperNoteUseCase creates a new use case instance
func NewUpdateWorkPaperNoteUseCase(deskService service.DeskService) *UpdateWorkPaperNoteUseCase {
	return &UpdateWorkPaperNoteUseCase{
		deskService: deskService,
	}
}

// UpdateWorkPaperNoteRequest represents the request payload for updating work paper note
type UpdateWorkPaperNoteRequest struct {
	ID         string `json:"id" validate:"required"`
	GDriveLink string `json:"gdrive_link,omitempty"`
	IsValid    *bool  `json:"is_valid,omitempty"`
	Notes      string `json:"notes,omitempty"`
}

// UpdateWorkPaperNoteResponse represents the response payload for updating work paper note
type UpdateWorkPaperNoteResponse struct {
	ID           string  `json:"id"`
	WorkPaperID  string  `json:"work_paper_id"`
	MasterItemID string  `json:"master_item_id"`
	GDriveLink   *string `json:"gdrive_link"`
	IsValid      *bool   `json:"is_valid"`
	Notes        *string `json:"notes"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

// Execute executes the use case
func (uc *UpdateWorkPaperNoteUseCase) Execute(ctx context.Context, req UpdateWorkPaperNoteRequest) (*UpdateWorkPaperNoteResponse, error) {
	// Get current work paper note
	currentNote, err := uc.deskService.GetWorkPaperNoteByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.GDriveLink != "" {
		_, err = uc.deskService.UpdateWorkPaperNoteLink(ctx, req.ID, req.GDriveLink)
		if err != nil {
			return nil, err
		}
	}

	// Update validation and notes if provided
	if req.IsValid != nil || req.Notes != "" {
		var isValid *bool
		notes := req.Notes

		// If only one is provided, keep the existing value for the other
		if req.IsValid == nil {
			isValid = currentNote.IsValid
		} else {
			isValid = req.IsValid
		}
		if notes == "" {
			notes = currentNote.GetNotes()
		}

		_, err = uc.deskService.UpdateWorkPaperNoteValidation(ctx, req.ID, isValid, notes)
		if err != nil {
			return nil, err
		}
	}

	// Get the updated note to return full response
	updatedNote, err := uc.deskService.GetWorkPaperNoteByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// Convert to response
	response := &UpdateWorkPaperNoteResponse{
		ID:           updatedNote.ID.String(),
		WorkPaperID:  updatedNote.WorkPaperID.String(),
		MasterItemID: updatedNote.MasterItemID.String(),
		GDriveLink:   updatedNote.GDriveLink,
		IsValid:      updatedNote.IsValid,
		Notes:        updatedNote.Notes,
		CreatedAt:    updatedNote.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    updatedNote.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return response, nil
}
