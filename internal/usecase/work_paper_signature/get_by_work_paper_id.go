package work_paper_signature

import (
	"context"

	"github.com/google/uuid"
	"sandbox/internal/domain/repository"
)

type GetWorkPaperSignaturesByWorkPaperIDUseCase struct {
	signatureRepo repository.WorkPaperSignatureRepository
}

func NewGetWorkPaperSignaturesByWorkPaperIDUseCase(signatureRepo repository.WorkPaperSignatureRepository) *GetWorkPaperSignaturesByWorkPaperIDUseCase {
	return &GetWorkPaperSignaturesByWorkPaperIDUseCase{
		signatureRepo: signatureRepo,
	}
}

func (uc *GetWorkPaperSignaturesByWorkPaperIDUseCase) Execute(ctx context.Context, workPaperID string) ([]*WorkPaperSignatureResponse, error) {
	// Parse work paper ID
	parsedWorkPaperID, err := uuid.Parse(workPaperID)
	if err != nil {
		return nil, ErrInvalidWorkPaperID
	}

	// Get signatures by work paper ID
	signatures, err := uc.signatureRepo.GetByWorkPaperID(ctx, parsedWorkPaperID)
	if err != nil {
		return nil, err
	}

	// Convert entities to response DTOs
	var responses []*WorkPaperSignatureResponse
	for _, signature := range signatures {
		responses = append(responses, FromEntity(signature))
	}

	return responses, nil
}
