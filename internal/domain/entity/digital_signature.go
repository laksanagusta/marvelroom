package entity

import (
	"time"
)

// DigitalSignature represents certificate-based signature data
type DigitalSignature struct {
	Signature         string     `json:"signature" db:"signature"`
	Payload           string     `json:"payload" db:"payload"`
	Algorithm         string     `json:"algorithm" db:"algorithm"`
	PublicKeyID       string     `json:"public_key_id" db:"public_key_id"`
	CertificateID     string     `json:"certificate_id" db:"certificate_id"`
	Timestamp         time.Time  `json:"timestamp" db:"timestamp"`
	Verified          bool       `json:"verified" db:"verified"`
	VerifiedAt        *time.Time `json:"verified_at" db:"verified_at"`
	VerificationError string     `json:"verification_error,omitempty" db:"verification_error"`
}

// DigitalSignatureVerification represents verification request/response
type DigitalSignatureVerification struct {
	Signature    string            `json:"signature" validate:"required"`
	Payload      *SignaturePayload `json:"payload" validate:"required"`
	UserID       string            `json:"user_id" validate:"required"`
	WorkPaperID  string            `json:"work_paper_id" validate:"required"`
	SignatureID  string            `json:"signature_id" validate:"required"`
	IsValid      bool              `json:"is_valid"`
	VerifiedAt   time.Time         `json:"verified_at"`
	Algorithm    string            `json:"algorithm"`
	ErrorMessage string            `json:"error_message,omitempty"`
}

// SignaturePayload represents the data structure that gets signed
type SignaturePayload struct {
	UserID               string    `json:"user_id" validate:"required"`
	WorkPaperID          string    `json:"work_paper_id" validate:"required"`
	WorkPaperSignatureID string    `json:"work_paper_signature_id" validate:"required"`
	Timestamp            time.Time `json:"timestamp" validate:"required"`
}

// NewDigitalSignature creates a new DigitalSignature entity
func NewDigitalSignature(signature, payload, algorithm, publicKeyID, certificateID string, timestamp time.Time) *DigitalSignature {
	return &DigitalSignature{
		Signature:     signature,
		Payload:       payload,
		Algorithm:     algorithm,
		PublicKeyID:   publicKeyID,
		CertificateID: certificateID,
		Timestamp:     timestamp,
		Verified:      false,
	}
}

// MarkVerified marks the signature as verified
func (ds *DigitalSignature) MarkVerified() {
	now := time.Now().UTC()
	ds.Verified = true
	ds.VerifiedAt = &now
	ds.VerificationError = ""
}

// MarkVerificationFailed marks the signature verification as failed
func (ds *DigitalSignature) MarkVerificationFailed(errorMsg string) {
	ds.Verified = false
	ds.VerifiedAt = nil
	ds.VerificationError = errorMsg
}

// IsValid checks if the signature is valid and verified
func (ds *DigitalSignature) IsValid() bool {
	return ds.Verified && ds.VerificationError == ""
}
