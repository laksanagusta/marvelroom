package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/pkg/database"
)

// SQL queries for business trip history
const (
	insertBusinessTripHistory = `
		INSERT INTO business_trip_histories (
			id, business_trip_id, change_type, field_name, old_value, new_value,
			changed_by, notes, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	findBusinessTripHistoryByID = `
		SELECT
			id, business_trip_id, change_type, field_name, old_value, new_value,
			changed_by, notes, created_at
		FROM business_trip_histories
		WHERE id = $1
	`

	findBusinessTripHistoriesByBusinessTripID = `
		SELECT
			id, business_trip_id, change_type, field_name, old_value, new_value,
			changed_by, notes, created_at
		FROM business_trip_histories
		WHERE business_trip_id = $1
		ORDER BY created_at DESC
	`

	findBusinessTripHistoriesByBusinessTripIDAndChangeType = `
		SELECT
			id, business_trip_id, change_type, field_name, old_value, new_value,
			changed_by, notes, created_at
		FROM business_trip_histories
		WHERE business_trip_id = $1 AND change_type = $2
		ORDER BY created_at DESC
	`
)

// NewBusinessTripHistoryRepository creates a new instance of BusinessTripHistoryRepository
func NewBusinessTripHistoryRepository(db database.Queryer) repository.BusinessTripHistoryRepository {
	return &businessTripHistoryRepository{
		db: db,
	}
}

type businessTripHistoryRepository struct {
	db database.Queryer
}

// Create creates a new business trip history record
func (r *businessTripHistoryRepository) Create(ctx context.Context, history *entity.BusinessTripHistory) error {
	if history.ID == "" {
		history.ID = uuid.New().String()
	}

	var returnedID string
	err := r.db.GetContext(ctx, &returnedID, insertBusinessTripHistory,
		history.ID,
		history.BusinessTripID,
		history.ChangeType,
		history.FieldName,
		history.OldValue,
		history.NewValue,
		history.ChangedBy,
		history.Notes,
		history.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create business trip history: %w", err)
	}

	if returnedID != history.ID {
		return fmt.Errorf("returned ID %s does not match expected ID %s", returnedID, history.ID)
	}

	return nil
}

// FindByBusinessTripID retrieves all history records for a specific business trip
func (r *businessTripHistoryRepository) FindByBusinessTripID(ctx context.Context, businessTripID string) ([]*entity.BusinessTripHistory, error) {
	rows, err := r.db.QueryxContext(ctx, findBusinessTripHistoriesByBusinessTripID, businessTripID)
	if err != nil {
		return nil, fmt.Errorf("failed to query business trip histories: %w", err)
	}
	defer rows.Close()

	var histories []*entity.BusinessTripHistory
	for rows.Next() {
		var history entity.BusinessTripHistory
		err := rows.StructScan(&history)
		if err != nil {
			return nil, fmt.Errorf("failed to scan business trip history: %w", err)
		}
		histories = append(histories, &history)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return histories, nil
}

// FindByID retrieves a history record by its ID
func (r *businessTripHistoryRepository) FindByID(ctx context.Context, id string) (*entity.BusinessTripHistory, error) {
	var history entity.BusinessTripHistory
	err := r.db.GetContext(ctx, &history, findBusinessTripHistoryByID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get business trip history: %w", err)
	}

	return &history, nil
}

// FindByBusinessTripIDAndChangeType retrieves history records by business trip ID and change type
func (r *businessTripHistoryRepository) FindByBusinessTripIDAndChangeType(ctx context.Context, businessTripID string, changeType entity.BusinessTripHistoryChangeType) ([]*entity.BusinessTripHistory, error) {
	rows, err := r.db.QueryxContext(ctx, findBusinessTripHistoriesByBusinessTripIDAndChangeType, businessTripID, changeType)
	if err != nil {
		return nil, fmt.Errorf("failed to query business trip histories by change type: %w", err)
	}
	defer rows.Close()

	var histories []*entity.BusinessTripHistory
	for rows.Next() {
		var history entity.BusinessTripHistory
		err := rows.StructScan(&history)
		if err != nil {
			return nil, fmt.Errorf("failed to scan business trip history: %w", err)
		}
		histories = append(histories, &history)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return histories, nil
}

// WithTransaction returns a new repository instance with the given transaction
func (r *businessTripHistoryRepository) WithTransaction(tx database.DBTx) repository.BusinessTripHistoryRepository {
	return &businessTripHistoryRepository{
		db: tx,
	}
}
