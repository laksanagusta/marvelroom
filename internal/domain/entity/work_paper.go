package entity

import (
	"time"

	"github.com/google/uuid"
)

// WorkPaper represents a working paper per organization per semester
type WorkPaper struct {
	ID             uuid.UUID      `db:"id"`
	OrganizationID uuid.UUID      `db:"organization_id"`
	Year           int            `db:"year"`
	Semester       int            `db:"semester"` // 1 or 2
	Status         string         `db:"status"`   // draft, ongoing, ready_to_sign, completed
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
	DeletedAt      *time.Time     `db:"deleted_at"`

	// Relations
	Organization *Organization        `db:"-"`
	Notes        []*WorkPaperNote     `db:"-"`
	Signatures   []*WorkPaperSignature `db:"-"`
}

// WorkPaperStatus constants
const (
	WorkPaperStatusDraft        = "draft"
	WorkPaperStatusOngoing     = "ongoing"
	WorkPaperStatusReadyToSign = "ready_to_sign"
	WorkPaperStatusCompleted   = "completed"
)

// NewWorkPaper creates a new work paper with validation
func NewWorkPaper(organizationID uuid.UUID, year, semester int) (*WorkPaper, error) {
	if organizationID == uuid.Nil {
		return nil, ErrOrganizationIDRequired
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
		Year:           year,
		Semester:       semester,
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
		WorkPaperStatusDraft:        {WorkPaperStatusDraft, WorkPaperStatusOngoing},
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