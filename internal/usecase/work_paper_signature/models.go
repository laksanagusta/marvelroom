package work_paper_signature

import (
	"errors"

	"sandbox/internal/domain/entity"
)

var (
	ErrInvalidWorkPaperID = errors.New("invalid work paper ID")
)

// WorkPaperSignatureResponse represents the response for work paper signature
type WorkPaperSignatureResponse struct {
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
	SignedAt      string                `json:"signed_at,omitempty"`
}
