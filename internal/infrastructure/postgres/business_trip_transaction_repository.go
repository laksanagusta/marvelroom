package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/pkg/database"
)

// SQL queries for transaction operations
const (
	getTransactionByIDQuery = `
		SELECT
			id, assignee_id, name, type, subtype, amount, total_night, subtotal, description, transport_detail, created_at, updated_at
		FROM assignee_transactions
		WHERE id = $1 AND deleted_at IS NULL
	`

	getTransactionsByAssigneeIDQuery = `
		SELECT
			id, assignee_id, name, type, subtype, amount, total_night, subtotal, description, transport_detail, created_at, updated_at
		FROM assignee_transactions
		WHERE assignee_id = $1 AND deleted_at IS NULL
		ORDER BY created_at
	`

	deleteTransactionsByAssigneeIDsQueryTemplate = `
		UPDATE assignee_transactions
		SET deleted_at = $1
		WHERE assignee_id IN (%s) AND deleted_at IS NULL
	`

	getTransactionTotalCountQuery = `
		SELECT COUNT(*)
		FROM assignee_transactions at
		INNER JOIN assignees a ON at.assignee_id = a.id
		INNER JOIN business_trips bt ON a.business_trip_id = bt.id
		WHERE at.deleted_at IS NULL AND a.deleted_at IS NULL AND bt.deleted_at IS NULL
	`

	getTransactionTotalCountWithDateRangeQuery = `
		SELECT COUNT(*)
		FROM assignee_transactions at
		INNER JOIN assignees a ON at.assignee_id = a.id
		INNER JOIN business_trips bt ON a.business_trip_id = bt.id
		WHERE at.deleted_at IS NULL AND a.deleted_at IS NULL AND bt.deleted_at IS NULL
		AND ($1::timestamp IS NULL OR bt.created_at >= $1)
		AND ($2::timestamp IS NULL OR bt.created_at <= $2)
	`

	getTypeStatsQuery = `
		SELECT
			at.type as transaction_type,
			COUNT(*) as average_amount,
			COALESCE(SUM(at.subtotal), 0) as total_amount
		FROM assignee_transactions at
		INNER JOIN assignees a ON at.assignee_id = a.id
		INNER JOIN business_trips bt ON a.business_trip_id = bt.id
		WHERE at.deleted_at IS NULL AND a.deleted_at IS NULL AND bt.deleted_at IS NULL
		AND ($1::timestamp IS NULL OR bt.created_at >= $1)
		AND ($2::timestamp IS NULL OR bt.created_at <= $2)
		GROUP BY at.type
		ORDER BY total_amount DESC
	`
)

// NewBusinessTripTransactionRepository creates a new instance of BusinessTripTransactionRepository
func NewBusinessTripTransactionRepository(db database.Queryer) repository.BusinessTripTransactionRepository {
	return &businessTripTransactionRepository{
		db: db,
	}
}

type businessTripTransactionRepository struct {
	db database.Queryer
}

// WithTransaction returns a new repository instance with given transaction
func (r *businessTripTransactionRepository) WithTransaction(tx database.DBTx) repository.BusinessTripTransactionRepository {
	return &businessTripTransactionRepository{
		db: tx,
	}
}

// CreateTransaction creates a new transaction
func (r *businessTripTransactionRepository) CreateTransaction(ctx context.Context, transaction *entity.Transaction) (*entity.Transaction, error) {
	if transaction.ID == "" {
		transaction.ID = uuid.New().String()
	}

	now := time.Now()

	// Use insert query from business_trip_repository.go
	query := `
		INSERT INTO assignee_transactions (
			id, assignee_id, name, type, subtype, amount, total_night, subtotal, description, transport_detail, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`

	var returnedID string
	err := r.db.GetContext(ctx, &returnedID, query,
		transaction.ID,
		transaction.AssigneeID,
		transaction.Name,
		transaction.Type,
		transaction.Subtype,
		transaction.Amount,
		transaction.TotalNight,
		transaction.Subtotal,
		transaction.Description,
		transaction.TransportDetail,
		now,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if returnedID != transaction.ID {
		return nil, fmt.Errorf("returned ID %s does not match expected ID %s", returnedID, transaction.ID)
	}

	// Set timestamps
	transaction.CreatedAt = now
	transaction.UpdatedAt = now

	return transaction, nil
}

// GetTransactionByID retrieves a transaction by ID
func (r *businessTripTransactionRepository) GetTransactionByID(ctx context.Context, id string) (*entity.Transaction, error) {
	var transaction entity.Transaction
	err := r.db.GetContext(ctx, &transaction, getTransactionByIDQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &transaction, nil
}

// UpdateTransaction updates a transaction
func (r *businessTripTransactionRepository) UpdateTransaction(ctx context.Context, transaction *entity.Transaction) (*entity.Transaction, error) {
	now := time.Now()

	// Use update query from business_trip_repository.go
	query := `
		UPDATE assignee_transactions
		SET name = $2, type = $3, subtype = $4, amount = $5, total_night = $6, subtotal = $7, description = $8, transport_detail = $9, updated_at = $10
		WHERE id = $1 AND deleted_at IS NULL
	`

	res, err := r.db.ExecContext(ctx, query,
		transaction.ID,
		transaction.Name,
		transaction.Type,
		transaction.Subtype,
		transaction.Amount,
		transaction.TotalNight,
		transaction.Subtotal,
		transaction.Description,
		transaction.TransportDetail,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowAffected == 0 {
		return nil, fmt.Errorf("transaction with ID %s not found", transaction.ID)
	}

	transaction.UpdatedAt = now

	return transaction, nil
}

// DeleteTransaction soft deletes a transaction
func (r *businessTripTransactionRepository) DeleteTransaction(ctx context.Context, id string) error {
	now := time.Now()

	// Use delete query from business_trip_repository.go
	query := `
		UPDATE assignee_transactions
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	res, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowAffected == 0 {
		return fmt.Errorf("transaction with ID %s not found", id)
	}

	return nil
}

// GetTransactionsByAssigneeID retrieves all transactions for an assignee
func (r *businessTripTransactionRepository) GetTransactionsByAssigneeID(ctx context.Context, assigneeID string) ([]*entity.Transaction, error) {
	rows, err := r.db.QueryxContext(ctx, getTransactionsByAssigneeIDQuery, assigneeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*entity.Transaction
	for rows.Next() {
		var transaction entity.Transaction
		err := rows.StructScan(&transaction)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, &transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return transactions, nil
}

// DeleteTransactionsByAssigneeIDs deletes transactions by multiple assignee IDs
func (r *businessTripTransactionRepository) DeleteTransactionsByAssigneeIDs(ctx context.Context, assigneeIDs []string) error {
	if len(assigneeIDs) == 0 {
		return nil
	}

	now := time.Now()

	// Build query dynamically for each assignee ID
	placeholders := make([]string, len(assigneeIDs))
	args := make([]interface{}, len(assigneeIDs)+1) // +1 for the timestamp

	args[0] = now // First argument is the timestamp
	for i, id := range assigneeIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+2) // Start from $2 since $1 is the timestamp
		args[i+1] = id
	}

	query := fmt.Sprintf(deleteTransactionsByAssigneeIDsQueryTemplate, strings.Join(placeholders, ","))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete transactions by assignee IDs: %w", err)
	}

	return nil
}

// GetTotalCount gets total count of transactions with optional date filtering
func (r *businessTripTransactionRepository) GetTotalCount(ctx context.Context, startDate, endDate *time.Time) (int64, error) {
	var count int64
	var err error

	if startDate != nil || endDate != nil {
		err = r.db.GetContext(ctx, &count, getTransactionTotalCountWithDateRangeQuery, startDate, endDate)
	} else {
		err = r.db.GetContext(ctx, &count, getTransactionTotalCountQuery)
	}

	if err != nil {
		return 0, fmt.Errorf("failed to get total transaction count: %w", err)
	}

	return count, nil
}

// GetTypeStats gets transaction type statistics for dashboard
func (r *businessTripTransactionRepository) GetTypeStats(ctx context.Context, startDate, endDate *time.Time) ([]*repository.TransactionTypeData, error) {
	rows, err := r.db.QueryxContext(ctx, getTypeStatsQuery, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query transaction type stats: %w", err)
	}
	defer rows.Close()

	var stats []*repository.TransactionTypeData
	for rows.Next() {
		var stat repository.TransactionTypeData
		err := rows.StructScan(&stat)
		if err != nil {
			log.Println("sdajda")
			return nil, fmt.Errorf("failed to scan transaction type stat: %w", err)
		}
		stats = append(stats, &stat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return stats, nil
}
