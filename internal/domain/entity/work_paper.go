package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// WorkPaper represents a working paper per organization per semester
type WorkPaper struct {
	ID               uuid.UUID  `db:"id"`
	OrganizationID   uuid.UUID  `db:"organization_id"`
	Name             string     `db:"name"`
	Year             int        `db:"year"`
	Semester         int        `db:"semester"` // 1 or 2
	TopicID          *uuid.UUID `db:"topic_id"`
	Status           string     `db:"status"`              // draft, ongoing, ready_to_sign, completed
	SourceFolderLink *string    `db:"source_folder_link"`  // Google Drive root folder link for smart linking
	LastFolderSyncAt *time.Time `db:"last_folder_sync_at"` // Last folder sync timestamp
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
	DeletedAt        *time.Time `db:"deleted_at"`

	// Relations
	Organization *Organization         `db:"-"`
	Notes        []*WorkPaperNote      `db:"-"`
	Signatures   []*WorkPaperSignature `db:"-"`
}

// WorkPaperStatus constants
const (
	WorkPaperStatusDraft       = "draft"
	WorkPaperStatusOngoing     = "ongoing"
	WorkPaperStatusReadyToSign = "ready_to_sign"
	WorkPaperStatusCompleted   = "completed"
)

// NewWorkPaper creates a new work paper with validation
func NewWorkPaper(organizationID uuid.UUID, name string, year, semester int, topicID *uuid.UUID) (*WorkPaper, error) {
	if organizationID == uuid.Nil {
		return nil, ErrOrganizationIDRequired
	}

	if name == "" {
		name = fmt.Sprintf("Work Paper %d S%d", year, semester)
	}

	if year < 2000 || year > 2100 {
		return nil, ErrInvalidYear
	}

	if semester != 1 && semester != 2 {
		return nil, ErrInvalidSemester
	}

	now := time.Now()
	return &WorkPaper{
		ID:             uuid.New(),
		OrganizationID: organizationID,
		Name:           name,
		Year:           year,
		Semester:       semester,
		TopicID:        topicID,
		Status:         WorkPaperStatusDraft,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// UpdateStatus updates the work paper status with validation
func (wp *WorkPaper) UpdateStatus(newStatus string) error {
	if !isValidStatusTransition(wp.Status, newStatus) {
		return ErrInvalidStatusTransition
	}

	wp.Status = newStatus
	wp.UpdatedAt = time.Now()
	return nil
}

// isValidStatusTransition validates if status transition is allowed
func isValidStatusTransition(currentStatus, newStatus string) bool {
	validTransitions := map[string][]string{
		WorkPaperStatusDraft:       {WorkPaperStatusDraft, WorkPaperStatusOngoing},
		WorkPaperStatusOngoing:     {WorkPaperStatusOngoing, WorkPaperStatusReadyToSign, WorkPaperStatusDraft},
		WorkPaperStatusReadyToSign: {WorkPaperStatusReadyToSign, WorkPaperStatusCompleted, WorkPaperStatusOngoing},
		WorkPaperStatusCompleted:   {WorkPaperStatusCompleted, WorkPaperStatusReadyToSign}, // Allow reopening if needed
	}

	allowedStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return false
	}

	for _, status := range allowedStatuses {
		if status == newStatus {
			return true
		}
	}

	return false
}
