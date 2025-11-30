package work_paper_signature

import (
	"context"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/pkg/pagination"
)

type ListWorkPaperSignaturesUseCase struct {
	signatureRepo repository.WorkPaperSignatureRepository
}

func NewListWorkPaperSignaturesUseCase(signatureRepo repository.WorkPaperSignatureRepository) *ListWorkPaperSignaturesUseCase {
	return &ListWorkPaperSignaturesUseCase{
		signatureRepo: signatureRepo,
	}
}

func (uc *ListWorkPaperSignaturesUseCase) Execute(ctx context.Context, params *pagination.QueryParams) ([]*WorkPaperSignatureResponse, *pagination.PagedResponse, error) {
	signatures, totalCount, err := uc.signatureRepo.List(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	// Convert entities to response DTOs
	var responses []*WorkPaperSignatureResponse
	for _, signature := range signatures {
		responses = append(responses, FromEntity(signature))
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

func FromEntity(signature *entity.WorkPaperSignature) *WorkPaperSignatureResponse {
	response := &WorkPaperSignatureResponse{
		ID:            signature.ID.String(),
		WorkPaperID:   signature.WorkPaperID.String(),
		UserID:        signature.UserID,
		UserName:      signature.UserName,
		UserEmail:     signature.GetUserEmail(),
		UserRole:      signature.GetUserRole(),
		SignatureType: signature.SignatureType,
		Status:        signature.Status,
		CreatedAt:     signature.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     signature.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	// Add signature data if available
	if signature.SignatureData != nil {
		response.SignatureData = signature.SignatureData
	}

	// Add signed timestamp if available
	if signature.SignedAt != nil {
		response.SignedAt = signature.SignedAt.Format("2006-01-02T15:04:05Z")
	}

	// Add notes if available
	response.Notes = signature.Notes

	return response
}
