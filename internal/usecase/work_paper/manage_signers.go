package work_paper

import (
	"context"

	"sandbox/internal/domain/service"
)

// ManageSignersUseCase handles signer management operations
type ManageSignersUseCase struct {
	deskService service.DeskService
}

// NewManageSignersUseCase creates a new use case instance
func NewManageSignersUseCase(deskService service.DeskService) *ManageSignersUseCase {
	return &ManageSignersUseCase{
		deskService: deskService,
	}
}

// Use types from service package instead

// Execute executes the use case for managing signers
func (uc *ManageSignersUseCase) Execute(ctx context.Context, req service.ManageSignersRequest) (*service.ManageSignersResponse, error) {
	switch req.Action {
	case "add":
		return uc.addSigners(ctx, req)
	case "remove":
		return uc.removeSigners(ctx, req)
	case "replace":
		return uc.replaceSigners(ctx, req)
	default:
		return nil, nil // This shouldn't happen due to validation
	}
}

// addSigners adds new signers to the work paper
func (uc *ManageSignersUseCase) addSigners(ctx context.Context, req service.ManageSignersRequest) (*service.ManageSignersResponse, error) {
	var signerResponses []service.SignerResponse

	for _, signer := range req.Signers {
		signatureReq := &service.CreateWorkPaperSignatureRequest{
			WorkPaperID:   req.WorkPaperID,
			UserID:        signer.UserID,
			UserName:      signer.UserName,
			UserEmail:     signer.UserEmail,
			UserRole:      signer.UserRole,
			SignatureType: signer.SignatureType,
		}

		signature, err := uc.deskService.CreateWorkPaperSignature(ctx, signatureReq)
		if err != nil {
			continue // Skip failed signers, continue with others
		}

		signerResponse := service.SignerResponse{
			SignatureID:   signature.ID.String(),
			UserID:        signature.UserID,
			UserName:      signature.UserName,
			UserEmail:     signature.GetUserEmail(),
			UserRole:      signature.GetUserRole(),
			SignatureType: signature.SignatureType,
			Status:        signature.Status,
			CreatedAt:     signature.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		signerResponses = append(signerResponses, signerResponse)
	}

	return &service.ManageSignersResponse{
		WorkPaperID: req.WorkPaperID,
		Action:      req.Action,
		Signers:     signerResponses,
		Message:     "Signers added successfully",
	}, nil
}

// removeSigners removes signers from the work paper
func (uc *ManageSignersUseCase) removeSigners(ctx context.Context, req service.ManageSignersRequest) (*service.ManageSignersResponse, error) {
	// For remove action, we only need user IDs from signers
	var signerResponses []service.SignerResponse

	for _, signer := range req.Signers {
		// Get existing signatures for this work paper and user
		signatures, err := uc.deskService.GetWorkPaperSignaturesByUserID(ctx, signer.UserID)
		if err != nil {
			continue
		}

		// Find and remove signatures for this work paper
		for _, sig := range signatures {
			if sig.WorkPaperID.String() == req.WorkPaperID {
				// Only remove pending signatures
				if sig.Status == "pending" {
					// Note: You might want to add a soft delete method in service
					// For now, we'll just mark it as rejected to remove it from active flow
					rejectReq := &service.RejectWorkPaperSignatureRequest{
						Notes: "Removed from signers",
					}

					_, err := uc.deskService.RejectWorkPaperSignature(ctx, sig.ID.String(), rejectReq)
					if err != nil {
						continue
					}

					signerResponse := service.SignerResponse{
						SignatureID: sig.ID.String(),
						UserID:      sig.UserID,
						UserName:    sig.UserName,
						Status:      "removed",
					}
					signerResponses = append(signerResponses, signerResponse)
				}
			}
		}
	}

	return &service.ManageSignersResponse{
		WorkPaperID: req.WorkPaperID,
		Action:      req.Action,
		Signers:     signerResponses,
		Message:     "Signers removed successfully",
	}, nil
}

// replaceSigners replaces all signers with new ones
func (uc *ManageSignersUseCase) replaceSigners(ctx context.Context, req service.ManageSignersRequest) (*service.ManageSignersResponse, error) {
	// First, get all existing signatures for this work paper
	existingSignatures, err := uc.deskService.GetWorkPaperSignatures(ctx, req.WorkPaperID)
	if err != nil {
		return nil, err
	}

	// Remove all existing pending signatures
	for _, sig := range existingSignatures {
		if sig.Status == "pending" {
			rejectReq := &service.RejectWorkPaperSignatureRequest{
				Notes: "Replaced with new signers",
			}
			uc.deskService.RejectWorkPaperSignature(ctx, sig.ID.String(), rejectReq)
		}
	}

	// Add new signers
	return uc.addSigners(ctx, req)
}
