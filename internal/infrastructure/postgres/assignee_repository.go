package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/pkg/database"
)

// SQL queries for assignee operations
const (
	getAssigneeByIDQuery = `
		SELECT
			id, business_trip_id, name, spd_number, employee_id, position, rank, employee_name, employee_number, created_at, updated_at
		FROM assignees
		WHERE id = $1 AND deleted_at IS NULL
	`

	getAssigneesByBusinessTripIDQuery = `
		SELECT
			id, business_trip_id, name, spd_number, employee_id, position, rank, employee_name, employee_number, created_at, updated_at
		FROM assignees
		WHERE business_trip_id = $1 AND deleted_at IS NULL
		ORDER BY created_at
	`

	getAssigneesByBusinessTripIDWithoutTransactionsQuery = `
		SELECT
			id, business_trip_id, name, spd_number, employee_id, position, rank, employee_name, employee_number, created_at, updated_at
		FROM assignees
		WHERE business_trip_id = $1 AND deleted_at IS NULL
		ORDER BY created_at
	`

	deleteAssigneesByBusinessTripIDQuery = `
		UPDATE assignees
		SET deleted_at = $1
		WHERE business_trip_id = $2 AND deleted_at IS NULL
	`

	hardDeleteAssigneesByBusinessTripIDQuery = `
		DELETE FROM assignees
		WHERE business_trip_id = $1
	`

	getAssigneeTotalCountQuery = `
		SELECT COUNT(*)
		FROM assignees a
		INNER JOIN business_trips bt ON a.business_trip_id = bt.id
		WHERE a.deleted_at IS NULL AND bt.deleted_at IS NULL
	`

	getAssigneeTotalCountWithDateRangeQuery = `
		SELECT COUNT(*)
		FROM assignees a
		INNER JOIN business_trips bt ON a.business_trip_id = bt.id
		WHERE a.deleted_at IS NULL AND bt.deleted_at IS NULL
		AND ($1::timestamp IS NULL OR bt.created_at >= $1)
		AND ($2::timestamp IS NULL OR bt.created_at <= $2)
	`
)

// NewAssigneeRepository creates a new instance of AssigneeRepository
func NewAssigneeRepository(db database.Queryer) repository.AssigneeRepository {
	return &assigneeRepository{
		db: db,
	}
}

type assigneeRepository struct {
	db database.Queryer
}

// WithTransaction returns a new repository instance with the given transaction
func (r *assigneeRepository) WithTransaction(tx database.DBTx) repository.AssigneeRepository {
	return &assigneeRepository{
		db: tx,
	}
}

// Create creates a new assignee
func (r *assigneeRepository) Create(ctx context.Context, assignee *entity.Assignee) (*entity.Assignee, error) {
	if assignee.ID == "" {
		assignee.ID = uuid.New().String()
	}

	now := time.Now()

	// Query based on actual table structure
	query := `
		INSERT INTO assignees (
			id, business_trip_id, name, spd_number, employee_id, position, rank, employee_name, employee_number, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`

	var returnedID string
	err := r.db.GetContext(ctx, &returnedID, query,
		assignee.ID,
		assignee.BusinessTripID,
		assignee.Name,
		assignee.SPDNumber,
		assignee.EmployeeID,
		assignee.Position,
		assignee.Rank,
		assignee.EmployeeName,
		assignee.EmployeeNumber,
		now,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create assignee: %w", err)
	}

	if returnedID != assignee.ID {
		return nil, fmt.Errorf("returned ID %s does not match expected ID %s", returnedID, assignee.ID)
	}

	// Set timestamps
	assignee.CreatedAt = now
	assignee.UpdatedAt = now

	return assignee, nil
}

// GetAssigneeByID retrieves an assignee by ID
func (r *assigneeRepository) GetAssigneeByID(ctx context.Context, id string) (*entity.Assignee, error) {
	var assignee entity.Assignee
	err := r.db.GetContext(ctx, &assignee, getAssigneeByIDQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get assignee: %w", err)
	}

	// Initialize empty transactions slice - transactions are loaded separately
	assignee.Transactions = make([]*entity.Transaction, 0)

	return &assignee, nil
}

// UpdateAssignee updates an assignee
func (r *assigneeRepository) UpdateAssignee(ctx context.Context, assignee *entity.Assignee) (*entity.Assignee, error) {
	now := time.Now()

	query := `
		UPDATE assignees
		SET name = $2, spd_number = $3, employee_id = $4, position = $5, rank = $6, employee_name = $7, employee_number = $8, updated_at = $9
		WHERE id = $1 AND deleted_at IS NULL
	`

	res, err := r.db.ExecContext(ctx, query,
		assignee.ID,
		assignee.Name,
		assignee.SPDNumber,
		assignee.EmployeeID,
		assignee.Position,
		assignee.Rank,
		assignee.EmployeeName,
		assignee.EmployeeNumber,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update assignee: %w", err)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowAffected == 0 {
		return nil, fmt.Errorf("assignee with ID %s not found", assignee.ID)
	}

	assignee.UpdatedAt = now

	return assignee, nil
}

// DeleteAssignee soft deletes an assignee
func (r *assigneeRepository) DeleteAssignee(ctx context.Context, id string) error {
	now := time.Now()

	query := `
		UPDATE assignees
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	res, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("failed to delete assignee: %w", err)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowAffected == 0 {
		return fmt.Errorf("assignee with ID %s not found", id)
	}

	return nil
}

// GetAssigneesByBusinessTripID retrieves all assignees for a business trip
func (r *assigneeRepository) GetAssigneesByBusinessTripID(ctx context.Context, businessTripID string) ([]*entity.Assignee, error) {
	rows, err := r.db.QueryxContext(ctx, getAssigneesByBusinessTripIDQuery, businessTripID)
	if err != nil {
		return nil, fmt.Errorf("failed to query assignees: %w", err)
	}
	defer rows.Close()

	var assignees []*entity.Assignee
	for rows.Next() {
		var assignee entity.Assignee
		err := rows.StructScan(&assignee)
		if err != nil {
			return nil, fmt.Errorf("failed to scan assignee: %w", err)
		}

		// Initialize empty transactions slice - transactions are loaded separately
		assignee.Transactions = make([]*entity.Transaction, 0)
		assignees = append(assignees, &assignee)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return assignees, nil
}

// GetAssigneesByBusinessTripIDWithoutTransactions retrieves all assignees for a business trip without loading their transactions
func (r *assigneeRepository) GetAssigneesByBusinessTripIDWithoutTransactions(ctx context.Context, businessTripID string) ([]*entity.Assignee, error) {
	rows, err := r.db.QueryxContext(ctx, getAssigneesByBusinessTripIDWithoutTransactionsQuery, businessTripID)
	if err != nil {
		return nil, fmt.Errorf("failed to query assignees: %w", err)
	}
	defer rows.Close()

	var assignees []*entity.Assignee
	for rows.Next() {
		var assignee entity.Assignee
		err := rows.StructScan(&assignee)
		if err != nil {
			return nil, fmt.Errorf("failed to scan assignee: %w", err)
		}

		// Initialize empty transactions slice - don't load transactions to avoid transaction issues
		assignee.Transactions = make([]*entity.Transaction, 0)
		assignees = append(assignees, &assignee)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return assignees, nil
}

// DeleteAssigneesByBusinessTripID deletes all assignees for a business trip
func (r *assigneeRepository) DeleteAssigneesByBusinessTripID(ctx context.Context, businessTripID string) error {
	// Use hard delete to avoid unique constraint violations when reinserting assignees
	// with the same spd_number for the same business_trip_id
	res, err := r.db.ExecContext(ctx, hardDeleteAssigneesByBusinessTripIDQuery, businessTripID)
	if err != nil {
		return fmt.Errorf("failed to hard delete assignees by business trip ID: %w", err)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowAffected == 0 {
		// Not an error if no assignees exist
		return nil
	}

	return nil
}

// GetTotalCount gets the total count of assignees with optional date filtering
func (r *assigneeRepository) GetTotalCount(ctx context.Context, startDate, endDate *time.Time) (int64, error) {
	var count int64
	var err error

	if startDate != nil || endDate != nil {
		err = r.db.GetContext(ctx, &count, getAssigneeTotalCountWithDateRangeQuery, startDate, endDate)
	} else {
		err = r.db.GetContext(ctx, &count, getAssigneeTotalCountQuery)
	}

	if err != nil {
		return 0, fmt.Errorf("failed to get total assignee count: %w", err)
	}

	return count, nil
}