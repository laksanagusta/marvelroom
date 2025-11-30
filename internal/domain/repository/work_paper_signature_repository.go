package repository

import (
	"context"
	"sandbox/internal/domain/entity"
	"sandbox/pkg/pagination"
	"time"

	"github.com/google/uuid"
)

// WorkPaperSignatureRepository defines the interface for work paper signature operations
type WorkPaperSignatureRepository interface {
	// Create creates a new work paper signature
	Create(ctx context.Context, signature *entity.WorkPaperSignature) error

	// GetByID gets a work paper signature by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.WorkPaperSignature, error)

	// GetByWorkPaperID gets all signatures for a work paper
	GetByWorkPaperID(ctx context.Context, workPaperID uuid.UUID) ([]*entity.WorkPaperSignature, error)

	// GetByWorkPaperIDAndUserID gets signature for a specific work paper and user
	GetByWorkPaperIDAndUserID(ctx context.Context, workPaperID uuid.UUID, userID string) (*entity.WorkPaperSignature, error)

	// GetByUserID gets all signatures by user ID
	GetByUserID(ctx context.Context, userID string) ([]*entity.WorkPaperSignature, error)

	// GetPendingByUserID gets all pending signatures by user ID
	GetPendingByUserID(ctx context.Context, userID string) ([]*entity.WorkPaperSignature, error)

	// GetPendingSignatures gets all pending signatures for a work paper
	GetPendingSignatures(ctx context.Context, workPaperID uuid.UUID) ([]*entity.WorkPaperSignature, error)

	// GetSignedSignatures gets all signed signatures for a work paper
	GetSignedSignatures(ctx context.Context, workPaperID uuid.UUID) ([]*entity.WorkPaperSignature, error)

	// List gets work paper signatures with filtering and pagination
	List(ctx context.Context, params *pagination.QueryParams) ([]*entity.WorkPaperSignature, int64, error)

	// Update updates a work paper signature
	Update(ctx context.Context, signature *entity.WorkPaperSignature) error

	// Delete soft deletes a work paper signature
	Delete(ctx context.Context, id uuid.UUID) error

	// GetSignaturesByStatus gets signatures by status for a work paper
	GetSignaturesByStatus(ctx context.Context, workPaperID uuid.UUID, status string) ([]*entity.WorkPaperSignature, error)

	// GetSignatureStats gets signature statistics for a work paper
	GetSignatureStats(ctx context.Context, workPaperID uuid.UUID) (*SignatureStats, error)

	// GetRecentSignatures gets recent signatures within a date range
	GetRecentSignatures(ctx context.Context, workPaperID uuid.UUID, from, to time.Time) ([]*entity.WorkPaperSignature, error)
}

// SignatureStats represents signature statistics
type SignatureStats struct {
	Total    int `json:"total"`
	Pending  int `json:"pending"`
	Signed   int `json:"signed"`
	Rejected int `json:"rejected"`
}
