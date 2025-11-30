package work_paper

import (
	"context"

	"sandbox/internal/domain/service"
)

// CheckWorkPaperNoteUseCase handles document checking using LLM
type CheckWorkPaperNoteUseCase struct {
	deskService service.DeskService
}

// NewCheckWorkPaperNoteUseCase creates a new use case instance
func NewCheckWorkPaperNoteUseCase(deskService service.DeskService) *CheckWorkPaperNoteUseCase {
	return &CheckWorkPaperNoteUseCase{
		deskService: deskService,
	}
}

// CheckRequest represents the request payload for checking a document
type CheckRequest struct {
	NoteID string `json:"note_id" validate:"required"`
}

// CheckResponse represents the response payload for checking a document
type CheckResponse struct {
	IsValid bool   `json:"is_valid"`
	Notes   string `json:"notes"`
	Model   string `json:"model"`
}

// Execute executes the use case
func (uc *CheckWorkPaperNoteUseCase) Execute(ctx context.Context, req CheckRequest) (*CheckResponse, error) {
	// Call service to check document
	checkResp, err := uc.deskService.CheckDocument(ctx, req.NoteID)
	if err != nil {
		return nil, err
	}

	// Convert to response
	response := &CheckResponse{
		IsValid: checkResp.IsValid,
		Notes:   checkResp.Notes,
		Model:   checkResp.Model,
	}

	return response, nil
}

// Backward compatibility aliases (deprecated)
type (
	CheckDocumentUseCase = CheckWorkPaperNoteUseCase
	CheckRequestLegacy   = CheckRequest
	CheckResponseLegacy  = CheckResponse
)

// NewCheckDocumentUseCase creates a new use case instance (deprecated)
func NewCheckDocumentUseCase(deskService service.DeskService) *CheckDocumentUseCase {
	return NewCheckWorkPaperNoteUseCase(deskService)
}
