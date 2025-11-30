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
	ErrWorkPaperSignatureNotFound  = errors.New("work paper signature not found")
	ErrCannotSignRejectedSignature = errors.New("cannot sign a rejected signature")
)

// CreateDigitalSignatureRequest represents the request for creating a digital signature
type CreateDigitalSignatureRequest struct {
	WorkPaperSignatureID string `json:"work_paper_signature_id" validate:"required"`
	UserID               string `json:"user_id" validate:"required"`
	Notes                string `json:"notes,omitempty"`
}

// CreateDigitalSignatureUseCase creates digital signatures
type CreateDigitalSignatureUseCase struct {
	workPaperSignatureRepo repository.WorkPaperSignatureRepository
	cryptoService          *cryptography.DigitalSignatureService
}

// NewCreateDigitalSignatureUseCase creates a new instance of CreateDigitalSignatureUseCase
func NewCreateDigitalSignatureUseCase(
	workPaperSignatureRepo repository.WorkPaperSignatureRepository,
	cryptoService *cryptography.DigitalSignatureService,
) *CreateDigitalSignatureUseCase {
	return &CreateDigitalSignatureUseCase{
		workPaperSignatureRepo: workPaperSignatureRepo,
		cryptoService:          cryptoService,
	}
}

// Execute creates a digital signature for a work paper signature
func (uc *CreateDigitalSignatureUseCase) Execute(request *CreateDigitalSignatureRequest) (*CreateDigitalSignatureResponse, error) {
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
			return nil, ErrWorkPaperSignatureNotFound
		}
		return nil, err
	}

	// Validate signature can be signed
	if signature.IsRejected() {
		return nil, ErrCannotSignRejectedSignature
	}

	if signature.IsSigned() {
		return nil, entity.ErrAlreadySigned
	}

	// Create payload for digital signature
	payload := cryptography.CreatePayloadFromData(
		request.UserID,
		signature.WorkPaperID.String(),
		signature.ID.String(),
	)

	// Create digital signature
	signatureResult, err := uc.cryptoService.SignPayload(payload)
	if err != nil {
		return nil, err
	}

	// Create digital signature entity
	digitalSignature := entity.NewDigitalSignature(
		signatureResult.Signature,
		signatureResult.Payload,
		signatureResult.Algorithm,
		"default", // publicKeyID - you can customize this
		"default", // certificateID - you can customize this
		signatureResult.Timestamp,
	)

	// Mark signature as verified (we just created it, so it should be valid)
	digitalSignature.MarkVerified()

	// Sign the work paper signature with digital signature
	err = signature.SignWithDigitalSignature(digitalSignature, request.Notes)
	if err != nil {
		return nil, err
	}

	// Update the signature in database
	err = uc.workPaperSignatureRepo.Update(ctx, signature)
	if err != nil {
		return nil, err
	}

	// Get updated signature
	updatedSignature, err := uc.workPaperSignatureRepo.GetByID(ctx, signature.ID)
	if err != nil {
		return nil, err
	}

	return &CreateDigitalSignatureResponse{
		ID:            updatedSignature.ID.String(),
		WorkPaperID:   updatedSignature.WorkPaperID.String(),
		UserID:        updatedSignature.UserID,
		UserName:      updatedSignature.UserName,
		UserEmail:     updatedSignature.GetUserEmail(),
		UserRole:      updatedSignature.GetUserRole(),
		SignatureType: updatedSignature.SignatureType,
		Status:        updatedSignature.Status,
		SignatureData: func() *entity.SignatureData {
			data := updatedSignature.GetSignatureData()
			if data != (entity.SignatureData{}) {
				return &data
			}
			return nil
		}(),
		Notes: func() *string {
			n := updatedSignature.GetNotes()
			if n != "" {
				return &n
			} else {
				return nil
			}
		}(),
		CreatedAt: updatedSignature.CreatedAt.Format(time.RFC3339),
		UpdatedAt: updatedSignature.UpdatedAt.Format(time.RFC3339),
		SignedAt: func() *string {
			if updatedSignature.SignedAt != nil {
				s := updatedSignature.SignedAt.Format(time.RFC3339)
				return &s
			} else {
				return nil
			}
		}(),
	}, nil
}

// CreateDigitalSignatureResponse represents the response after creating a digital signature
type CreateDigitalSignatureResponse struct {
	ID            string                `json:"id"`
	WorkPaperID   string                `json:"work_paper_id"`
	UserID        string                `json:"user_id"`
	UserName      string                `json:"user_name"`
	UserEmail     string                `json:"user_email,omitempty"`
	UserRole      string                `json:"user_role,omitempty"`
	SignatureType string                `json:"signature_type"`
	Status        string                `json:"status"`
	SignatureData *entity.SignatureData `json:"signature_data,omitempty"`
	Notes         *string               `json:"notes,omitempty"`
	CreatedAt     string                `json:"created_at"`
	UpdatedAt     string                `json:"updated_at"`
	SignedAt      *string               `json:"signed_at,omitempty"`
}

// DigitalSignatureResponse represents digital signature data in responses
type DigitalSignatureResponse struct {
	Signature   string  `json:"signature"`
	Payload     string  `json:"payload"`
	Algorithm   string  `json:"algorithm"`
	PublicKeyID string  `json:"public_key_id"`
	Timestamp   string  `json:"timestamp"`
	Verified    bool    `json:"verified"`
	VerifiedAt  *string `json:"verified_at,omitempty"`
}
