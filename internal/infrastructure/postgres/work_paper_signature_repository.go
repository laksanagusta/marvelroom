package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/pkg/pagination"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type workPaperSignatureRepository struct {
	db *sqlx.DB
}

// NewWorkPaperSignatureRepository creates a new work paper signature repository
func NewWorkPaperSignatureRepository(db *sqlx.DB) repository.WorkPaperSignatureRepository {
	return &workPaperSignatureRepository{
		db: db,
	}
}

// Create creates a new work paper signature
func (r *workPaperSignatureRepository) Create(ctx context.Context, signature *entity.WorkPaperSignature) error {
	query := `
		INSERT INTO work_paper_signatures (
			id, work_paper_id, user_id, user_name, user_email, user_role,
			signature_data, signature_type, status, notes, created_at, updated_at
		) VALUES (
			:id, :work_paper_id, :user_id, :user_name, :user_email, :user_role,
			:signature_data, :signature_type, :status, :notes, :created_at, :updated_at
		)`

	_, err := r.db.NamedExecContext(ctx, query, signature)
	if err != nil {
		return fmt.Errorf("failed to create work paper signature: %w", err)
	}

	return nil
}

// GetByID gets a work paper signature by ID
func (r *workPaperSignatureRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.WorkPaperSignature, error) {
	query := `
		SELECT id, work_paper_id, user_id, user_name, user_email, user_role,
			   signature_data, signed_at, signature_type, status, notes, created_at, updated_at, deleted_at
		FROM work_paper_signatures
		WHERE id = $1 AND deleted_at IS NULL`

	var signature entity.WorkPaperSignature
	err := r.db.GetContext(ctx, &signature, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrSignatureNotFound
		}
		return nil, fmt.Errorf("failed to get work paper signature by ID: %w", err)
	}

	return &signature, nil
}

// GetByWorkPaperID gets all signatures for a work paper
func (r *workPaperSignatureRepository) GetByWorkPaperID(ctx context.Context, workPaperID uuid.UUID) ([]*entity.WorkPaperSignature, error) {
	query := `
		SELECT id, work_paper_id, user_id, user_name, user_email, user_role,
			   signature_data, signed_at, signature_type, status, notes, created_at, updated_at, deleted_at
		FROM work_paper_signatures
		WHERE work_paper_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC`

	var signatures []*entity.WorkPaperSignature
	err := r.db.SelectContext(ctx, &signatures, query, workPaperID)
	if err != nil {
		return nil, fmt.Errorf("failed to get signatures by work paper ID: %w", err)
	}

	return signatures, nil
}

// GetByWorkPaperIDAndUserID gets signature for a specific work paper and user
func (r *workPaperSignatureRepository) GetByWorkPaperIDAndUserID(ctx context.Context, workPaperID uuid.UUID, userID string) (*entity.WorkPaperSignature, error) {
	query := `
		SELECT id, work_paper_id, user_id, user_name, user_email, user_role,
			   signature_data, signed_at, signature_type, status, notes, created_at, updated_at, deleted_at
		FROM work_paper_signatures
		WHERE work_paper_id = $1 AND user_id = $2 AND deleted_at IS NULL`

	var signature entity.WorkPaperSignature
	err := r.db.GetContext(ctx, &signature, query, workPaperID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrSignatureNotFound
		}
		return nil, fmt.Errorf("failed to get signature by work paper ID and user ID: %w", err)
	}

	return &signature, nil
}

// GetByUserID gets all signatures by user ID
func (r *workPaperSignatureRepository) GetByUserID(ctx context.Context, userID string) ([]*entity.WorkPaperSignature, error) {
	query := `
		SELECT id, work_paper_id, user_id, user_name, user_email, user_role,
			   signature_data, signed_at, signature_type, status, notes, created_at, updated_at, deleted_at
		FROM work_paper_signatures
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC`

	var signatures []*entity.WorkPaperSignature
	err := r.db.SelectContext(ctx, &signatures, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get signatures by user ID: %w", err)
	}

	return signatures, nil
}

// GetPendingByUserID gets all pending signatures by user ID
func (r *workPaperSignatureRepository) GetPendingByUserID(ctx context.Context, userID string) ([]*entity.WorkPaperSignature, error) {
	query := `
		SELECT id, work_paper_id, user_id, user_name, user_email, user_role,
			   signature_data, signed_at, signature_type, status, notes, created_at, updated_at, deleted_at
		FROM work_paper_signatures
		WHERE user_id = $1 AND status = $2 AND deleted_at IS NULL
		ORDER BY created_at ASC`

	var signatures []*entity.WorkPaperSignature
	err := r.db.SelectContext(ctx, &signatures, query, userID, entity.SignatureStatusPending)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending signatures by user ID: %w", err)
	}

	return signatures, nil
}

// GetPendingSignatures gets all pending signatures for a work paper
func (r *workPaperSignatureRepository) GetPendingSignatures(ctx context.Context, workPaperID uuid.UUID) ([]*entity.WorkPaperSignature, error) {
	return r.GetSignaturesByStatus(ctx, workPaperID, entity.SignatureStatusPending)
}

// GetSignedSignatures gets all signed signatures for a work paper
func (r *workPaperSignatureRepository) GetSignedSignatures(ctx context.Context, workPaperID uuid.UUID) ([]*entity.WorkPaperSignature, error) {
	return r.GetSignaturesByStatus(ctx, workPaperID, entity.SignatureStatusSigned)
}

// List gets work paper signatures with filtering and pagination
func (r *workPaperSignatureRepository) List(ctx context.Context, params *pagination.QueryParams) ([]*entity.WorkPaperSignature, int64, error) {
	// Build count query
	countBuilder := pagination.NewQueryBuilder("SELECT COUNT(*) FROM work_paper_signatures")
	for _, filter := range params.Filters {
		if err := countBuilder.AddFilter(filter); err != nil {
			return nil, 0, err
		}
	}

	// Always include deleted_at filter
	countBuilder.AddFilter(pagination.Filter{
		Field:    "deleted_at",
		Operator: "is",
		Value:    nil,
	})

	countQuery, countArgs := countBuilder.Build()

	var totalCount int64
	err := r.db.GetContext(ctx, &totalCount, countQuery, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	// Build main query
	queryBuilder := pagination.NewQueryBuilder(`
		SELECT id, work_paper_id, user_id, user_name, user_email, user_role,
			   signature_data, signed_at, signature_type, status, notes, created_at, updated_at, deleted_at
		FROM work_paper_signatures`)

	for _, filter := range params.Filters {
		if err := queryBuilder.AddFilter(filter); err != nil {
			return nil, 0, err
		}
	}

	// Always include deleted_at filter
	queryBuilder.AddFilter(pagination.Filter{
		Field:    "deleted_at",
		Operator: "is",
		Value:    nil,
	})

	for _, sort := range params.Sorts {
		if err := queryBuilder.AddSort(sort); err != nil {
			return nil, 0, err
		}
	}

	query, args := queryBuilder.Build()

	// Add pagination
	offset := (params.Pagination.Page - 1) * params.Pagination.Limit
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", params.Pagination.Limit, offset)

	var signatures []*entity.WorkPaperSignature
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var signature entity.WorkPaperSignature
		if err := rows.StructScan(&signature); err != nil {
			return nil, 0, fmt.Errorf("failed to scan work paper signature: %w", err)
		}
		signatures = append(signatures, &signature)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return signatures, totalCount, nil
}

// Update updates a work paper signature
func (r *workPaperSignatureRepository) Update(ctx context.Context, signature *entity.WorkPaperSignature) error {
	query := `
		UPDATE work_paper_signatures
		SET user_email = :user_email, user_role = :user_role, signature_data = :signature_data,
			signed_at = :signed_at, signature_type = :signature_type, status = :status,
			notes = :notes, updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL`

	_, err := r.db.NamedExecContext(ctx, query, signature)
	if err != nil {
		return fmt.Errorf("failed to update work paper signature: %w", err)
	}

	return nil
}

// Delete soft deletes a work paper signature
func (r *workPaperSignatureRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE work_paper_signatures SET deleted_at = $1 WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete work paper signature: %w", err)
	}

	return nil
}

// GetSignaturesByStatus gets signatures by status for a work paper
func (r *workPaperSignatureRepository) GetSignaturesByStatus(ctx context.Context, workPaperID uuid.UUID, status string) ([]*entity.WorkPaperSignature, error) {
	query := `
		SELECT id, work_paper_id, user_id, user_name, user_email, user_role,
			   signature_data, signed_at, signature_type, status, notes, created_at, updated_at, deleted_at
		FROM work_paper_signatures
		WHERE work_paper_id = $1 AND status = $2 AND deleted_at IS NULL
		ORDER BY created_at ASC`

	var signatures []*entity.WorkPaperSignature
	err := r.db.SelectContext(ctx, &signatures, query, workPaperID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get signatures by status: %w", err)
	}

	return signatures, nil
}

// GetSignatureStats gets signature statistics for a work paper
func (r *workPaperSignatureRepository) GetSignatureStats(ctx context.Context, workPaperID uuid.UUID) (*repository.SignatureStats, error) {
	query := `
		SELECT
			COUNT(*) as total,
			COUNT(CASE WHEN status = $1 THEN 1 END) as pending,
			COUNT(CASE WHEN status = $2 THEN 1 END) as signed,
			COUNT(CASE WHEN status = $3 THEN 1 END) as rejected
		FROM work_paper_signatures
		WHERE work_paper_id = $4 AND deleted_at IS NULL`

	var stats repository.SignatureStats
	err := r.db.GetContext(ctx, &stats, query,
		entity.SignatureStatusPending,
		entity.SignatureStatusSigned,
		entity.SignatureStatusRejected,
		workPaperID)
	if err != nil {
		return nil, fmt.Errorf("failed to get signature statistics: %w", err)
	}

	return &stats, nil
}

// GetRecentSignatures gets recent signatures within a date range
func (r *workPaperSignatureRepository) GetRecentSignatures(ctx context.Context, workPaperID uuid.UUID, from, to time.Time) ([]*entity.WorkPaperSignature, error) {
	query := `
		SELECT id, work_paper_id, user_id, user_name, user_email, user_role,
			   signature_data, signed_at, signature_type, status, notes, created_at, updated_at, deleted_at
		FROM work_paper_signatures
		WHERE work_paper_id = $1 AND created_at BETWEEN $2 AND $3 AND deleted_at IS NULL
		ORDER BY created_at DESC`

	var signatures []*entity.WorkPaperSignature
	err := r.db.SelectContext(ctx, &signatures, query, workPaperID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent signatures: %w", err)
	}

	return signatures, nil
}
