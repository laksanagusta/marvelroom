package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// WorkPaperSignature represents a signature on a work paper
type WorkPaperSignature struct {
	ID            uuid.UUID      `db:"id"`
	WorkPaperID   uuid.UUID      `db:"work_paper_id"`
	UserID        string         `db:"user_id"`
	UserName      string         `db:"user_name"`
	UserEmail     *string        `db:"user_email"`     // Nullable
	UserRole      *string        `db:"user_role"`      // Nullable
	SignatureData *SignatureData `db:"signature_data"` // Nullable JSONB
	SignedAt      *time.Time     `db:"signed_at"`      // Nullable
	SignatureType string         `db:"signature_type"`
	Status        string         `db:"status"`
	Notes         *string        `db:"notes"`           // Nullable
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at"`
	DeletedAt     *time.Time     `db:"deleted_at"`

	// Relations
	WorkPaper *WorkPaper `db:"-"`
}

// SignatureData represents the signature information
type SignatureData struct {
	SignatureImage string    `json:"signature_image,omitempty"` // Base64 encoded image
	IPAddress      string    `json:"ip_address,omitempty"`
	UserAgent      string    `json:"user_agent,omitempty"`
	DeviceID       string    `json:"device_id,omitempty"`
	Timestamp      time.Time `json:"timestamp"`
	Location       *Location `json:"location,omitempty"`
}

// Location represents geographic location
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address,omitempty"`
	City      string  `json:"city,omitempty"`
	Country   string  `json:"country,omitempty"`
}

// SignatureType constants
const (
	SignatureTypeDigital  = "digital"
	SignatureTypeManual   = "manual"
	SignatureTypeApproval = "approval"
)

// SignatureStatus constants
const (
	SignatureStatusPending = "pending"
	SignatureStatusSigned  = "signed"
	SignatureStatusRejected = "rejected"
)

// Value implements driver.Valuer interface for SignatureData
func (sd SignatureData) Value() (driver.Value, error) {
	return json.Marshal(sd)
}

// Scan implements sql.Scanner interface for SignatureData
func (sd *SignatureData) Scan(value interface{}) error {
	if value == nil {
		*sd = SignatureData{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, sd)
}

// NewWorkPaperSignature creates a new work paper signature
func NewWorkPaperSignature(workPaperID uuid.UUID, userID, userName, signatureType string) (*WorkPaperSignature, error) {
	if workPaperID == uuid.Nil {
		return nil, ErrWorkPaperIDRequired
	}

	if userID == "" {
		return nil, ErrUserIDRequired
	}

	if userName == "" {
		return nil, ErrUserNameRequired
	}

	if !isValidSignatureType(signatureType) {
		return nil, ErrInvalidSignatureType
	}

	now := time.Now()
	return &WorkPaperSignature{
		ID:           uuid.New(),
		WorkPaperID:  workPaperID,
		UserID:       userID,
		UserName:     userName,
		SignatureType: signatureType,
		Status:       SignatureStatusPending,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// SetUserDetails sets user details for the signature
func (wps *WorkPaperSignature) SetUserDetails(email, role string) {
	if email == "" {
		wps.UserEmail = nil
	} else {
		wps.UserEmail = &email
	}

	if role == "" {
		wps.UserRole = nil
	} else {
		wps.UserRole = &role
	}
	wps.UpdatedAt = time.Now()
}

// AddSignatureData adds signature data
func (wps *WorkPaperSignature) AddSignatureData(data SignatureData) {
	wps.SignatureData = &data
	wps.UpdatedAt = time.Now()
}

// Sign signs the work paper
func (wps *WorkPaperSignature) Sign(notes string) error {
	if wps.Status == SignatureStatusSigned {
		return ErrAlreadySigned
	}

	if wps.Status == SignatureStatusRejected {
		return ErrSignatureRejected
	}

	now := time.Now()
	wps.SignedAt = &now
	wps.Status = SignatureStatusSigned

	if notes == "" {
		wps.Notes = nil
	} else {
		wps.Notes = &notes
	}
	wps.UpdatedAt = time.Now()

	return nil
}

// Reject rejects the work paper signature
func (wps *WorkPaperSignature) Reject(notes string) error {
	if wps.Status == SignatureStatusSigned {
		return ErrAlreadySigned
	}

	wps.Status = SignatureStatusRejected
	if notes == "" {
		wps.Notes = nil
	} else {
		wps.Notes = &notes
	}
	wps.UpdatedAt = time.Now()

	return nil
}

// Reset resets the signature to pending status
func (wps *WorkPaperSignature) Reset() error {
	wps.Status = SignatureStatusPending
	wps.SignedAt = nil
	wps.Notes = nil
	wps.UpdatedAt = time.Now()
	return nil
}

// Helper methods to safely access nullable fields

// GetUserEmail safely returns the user email
func (wps *WorkPaperSignature) GetUserEmail() string {
	if wps.UserEmail == nil {
		return ""
	}
	return *wps.UserEmail
}

// GetUserRole safely returns the user role
func (wps *WorkPaperSignature) GetUserRole() string {
	if wps.UserRole == nil {
		return ""
	}
	return *wps.UserRole
}

// GetSignatureData safely returns the signature data
func (wps *WorkPaperSignature) GetSignatureData() SignatureData {
	if wps.SignatureData == nil {
		return SignatureData{}
	}
	return *wps.SignatureData
}

// GetNotes safely returns the notes
func (wps *WorkPaperSignature) GetNotes() string {
	if wps.Notes == nil {
		return ""
	}
	return *wps.Notes
}

// GetSignedAt safely returns the signed timestamp
func (wps *WorkPaperSignature) GetSignedAt() time.Time {
	if wps.SignedAt == nil {
		return time.Time{}
	}
	return *wps.SignedAt
}

// IsSigned returns true if the signature is signed
func (wps *WorkPaperSignature) IsSigned() bool {
	return wps.Status == SignatureStatusSigned
}

// IsPending returns true if the signature is pending
func (wps *WorkPaperSignature) IsPending() bool {
	return wps.Status == SignatureStatusPending
}

// IsRejected returns true if the signature is rejected
func (wps *WorkPaperSignature) IsRejected() bool {
	return wps.Status == SignatureStatusRejected
}

// isValidSignatureType validates if signature type is valid
func isValidSignatureType(signatureType string) bool {
	switch signatureType {
	case SignatureTypeDigital, SignatureTypeManual, SignatureTypeApproval:
		return true
	default:
		return false
	}
}