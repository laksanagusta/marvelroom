package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// BusinessTripHistoryChangeType represents the type of change in history
type BusinessTripHistoryChangeType string

const (
	HistoryChangeTypeStatusChange         BusinessTripHistoryChangeType = "status_change"
	HistoryChangeTypeVerificationApproved BusinessTripHistoryChangeType = "verification_approved"
	HistoryChangeTypeVerificationRejected BusinessTripHistoryChangeType = "verification_rejected"
	HistoryChangeTypeVerificationPending  BusinessTripHistoryChangeType = "verification_pending"
)

// BusinessTripHistory represents a history record of changes made to a business trip
type BusinessTripHistory struct {
	ID             string                        `db:"id"`
	BusinessTripID string                        `db:"business_trip_id"`
	ChangeType     BusinessTripHistoryChangeType `db:"change_type"`
	FieldName      sql.NullString                `db:"field_name"`
	OldValue       sql.NullString                `db:"old_value"`
	NewValue       sql.NullString                `db:"new_value"`
	ChangedBy      sql.NullString                `db:"changed_by"`
	Notes          sql.NullString                `db:"notes"`
	CreatedAt      time.Time                     `db:"created_at"`
}

// NewBusinessTripHistory creates a new business trip history record
func NewBusinessTripHistory(
	businessTripID string,
	changeType BusinessTripHistoryChangeType,
) (*BusinessTripHistory, error) {
	if businessTripID == "" {
		return nil, ErrInvalidInput
	}

	if changeType == "" {
		return nil, ErrInvalidInput
	}

	// Validate change type
	validChangeTypes := []BusinessTripHistoryChangeType{
		HistoryChangeTypeStatusChange,
		HistoryChangeTypeVerificationApproved,
		HistoryChangeTypeVerificationRejected,
		HistoryChangeTypeVerificationPending,
	}

	isValidType := false
	for _, validType := range validChangeTypes {
		if changeType == validType {
			isValidType = true
			break
		}
	}

	if !isValidType {
		return nil, ErrInvalidInput
	}

	return &BusinessTripHistory{
		ID:             uuid.New().String(),
		BusinessTripID: businessTripID,
		ChangeType:     changeType,
		CreatedAt:      time.Now(),
	}, nil
}

// SetFieldName sets the field name that was changed
func (h *BusinessTripHistory) SetFieldName(fieldName string) {
	if fieldName != "" {
		h.FieldName = sql.NullString{String: fieldName, Valid: true}
	}
}

// SetOldValue sets the old value before the change
func (h *BusinessTripHistory) SetOldValue(oldValue string) {
	if oldValue != "" {
		h.OldValue = sql.NullString{String: oldValue, Valid: true}
	}
}

// SetNewValue sets the new value after the change
func (h *BusinessTripHistory) SetNewValue(newValue string) {
	if newValue != "" {
		h.NewValue = sql.NullString{String: newValue, Valid: true}
	}
}

// SetChangedBy sets who made the change
func (h *BusinessTripHistory) SetChangedBy(username string) {
	if username != "" {
		h.ChangedBy = sql.NullString{String: username, Valid: true}
	}
}

// SetNotes sets additional notes about the change
func (h *BusinessTripHistory) SetNotes(notes string) {
	if notes != "" {
		h.Notes = sql.NullString{String: notes, Valid: true}
	}
}

// GetID returns the history ID
func (h *BusinessTripHistory) GetID() string {
	return h.ID
}

// GetBusinessTripID returns the business trip ID
func (h *BusinessTripHistory) GetBusinessTripID() string {
	return h.BusinessTripID
}

// GetChangeType returns the change type
func (h *BusinessTripHistory) GetChangeType() BusinessTripHistoryChangeType {
	return h.ChangeType
}

// GetFieldName returns the field name that was changed
func (h *BusinessTripHistory) GetFieldName() string {
	if h.FieldName.Valid {
		return h.FieldName.String
	}
	return ""
}

// GetOldValue returns the old value before the change
func (h *BusinessTripHistory) GetOldValue() string {
	if h.OldValue.Valid {
		return h.OldValue.String
	}
	return ""
}

// GetNewValue returns the new value after the change
func (h *BusinessTripHistory) GetNewValue() string {
	if h.NewValue.Valid {
		return h.NewValue.String
	}
	return ""
}

// GetChangedBy returns the username who made the change
func (h *BusinessTripHistory) GetChangedBy() string {
	if h.ChangedBy.Valid {
		return h.ChangedBy.String
	}
	return ""
}

// GetNotes returns the notes about the change
func (h *BusinessTripHistory) GetNotes() string {
	if h.Notes.Valid {
		return h.Notes.String
	}
	return ""
}

// GetCreatedAt returns the timestamp when the history was created
func (h *BusinessTripHistory) GetCreatedAt() time.Time {
	return h.CreatedAt
}
