package work_paper_signature

import (
	"context"
	"errors"
	"time"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/internal/infrastructure/cryptography"

	"github.com/google/uuid"
)

var (
	ErrSignatureNotFound  = errors.New("signature not found")
	ErrNoDigitalSignature = errors.New("no digital signature found")
)

// VerifyDigitalSignatureRequest represents the request for verifying a digital signature
type VerifyDigitalSignatureRequest struct {
	WorkPaperSignatureID string `json:"work_paper_signature_id" validate:"required"`
}

// VerifyDigitalSignatureUseCase verifies digital signatures
type VerifyDigitalSignatureUseCase struct {
	workPaperSignatureRepo repository.WorkPaperSignatureRepository
	cryptoService          *cryptography.DigitalSignatureService
}

// NewVerifyDigitalSignatureUseCase creates a new instance of VerifyDigitalSignatureUseCase
func NewVerifyDigitalSignatureUseCase(
	workPaperSignatureRepo repository.WorkPaperSignatureRepository,
	cryptoService *cryptography.DigitalSignatureService,
) *VerifyDigitalSignatureUseCase {
	return &VerifyDigitalSignatureUseCase{
		workPaperSignatureRepo: workPaperSignatureRepo,
		cryptoService:          cryptoService,
	}
}

// Execute verifies a digital signature for a work paper signature
func (uc *VerifyDigitalSignatureUseCase) Execute(request *VerifyDigitalSignatureRequest) (*VerifyDigitalSignatureResponse, error) {
	ctx := context.Background()

	// Parse work paper signature ID
	signatureID, err := uuid.Parse(request.WorkPaperSignatureID)
	if err != nil {
		return nil, err
	}

	// Get the work paper signature
	signature, err := uc.workPaperSignatureRepo.GetByID(ctx, signatureID)
	if err != nil {
		if errors.Is(err, entity.ErrSignatureNotFound) {
			return nil, ErrSignatureNotFound
		}
		return nil, err
	}

	// Get digital signature
	digitalSignature := signature.GetDigitalSignature()
	if digitalSignature == nil {
		return nil, ErrNoDigitalSignature
	}

	// Create payload for verification
	payload := &cryptography.SignaturePayload{
		UserID:               signature.UserID,
		WorkPaperID:          signature.WorkPaperID.String(),
		WorkPaperSignatureID: signature.ID.String(),
		Timestamp:            digitalSignature.Timestamp,
	}

	// Verify the signature
	err = uc.cryptoService.VerifySignature(digitalSignature.Signature, payload)
	if err != nil {
		// Mark signature as verification failed
		digitalSignature.MarkVerificationFailed(err.Error())

		// Update the signature in database
		updateErr := uc.workPaperSignatureRepo.Update(ctx, signature)
		if updateErr != nil {
			return nil, updateErr
		}

		return &VerifyDigitalSignatureResponse{
			WorkPaperSignatureID: request.WorkPaperSignatureID,
			IsValid:              false,
			VerifiedAt:           time.Now().Format(time.RFC3339),
			Algorithm:            digitalSignature.Algorithm,
			ErrorMessage:         err.Error(),
		}, nil
	}

	// Mark signature as verified
	digitalSignature.MarkVerified()

	// Update the signature in database
	err = uc.workPaperSignatureRepo.Update(ctx, signature)
	if err != nil {
		return nil, err
	}

	return &VerifyDigitalSignatureResponse{
		WorkPaperSignatureID: request.WorkPaperSignatureID,
		IsValid:              true,
		VerifiedAt:           time.Now().Format(time.RFC3339),
		Algorithm:            digitalSignature.Algorithm,
		ErrorMessage:         "",
	}, nil
}

// VerifyDigitalSignatureResponse represents the response after verifying a digital signature
type VerifyDigitalSignatureResponse struct {
	WorkPaperSignatureID string `json:"work_paper_signature_id"`
	IsValid              bool   `json:"is_valid"`
	VerifiedAt           string `json:"verified_at"`
	Algorithm            string `json:"algorithm"`
	ErrorMessage         string `json:"error_message,omitempty"`
}
