package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/pkg/business_trip_number"
	"sandbox/pkg/database"
	"sandbox/pkg/pagination"
)

// SQL queries
const (
	insertBusinessTrip = `
		INSERT INTO business_trips (
			id, business_trip_number, start_date, end_date, activity_purpose, destination_city,
			spd_date, departure_date, return_date, status, document_link, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`

	updateBusinessTrip = `
		UPDATE business_trips
		SET start_date = $2, end_date = $3, activity_purpose = $4, destination_city = $5,
			spd_date = $6, departure_date = $7, return_date = $8, status = $9, document_link = $10, updated_at = $11
		WHERE id = $1
	`

	findBusinessTripByID = `
		SELECT
			bt.id, bt.business_trip_number, bt.start_date, bt.end_date, bt.activity_purpose, bt.destination_city,
			bt.spd_date, bt.departure_date, bt.return_date, bt.status, bt.document_link, bt.created_at, bt.updated_at
		FROM business_trips bt
		WHERE bt.id = $1 AND bt.deleted_at IS NULL
	`

	deleteBusinessTrip = `
		UPDATE business_trips
		SET deleted_at = $1
		WHERE id = $2
	`

	insertAssignee = `
		INSERT INTO assignees (
			id, business_trip_id, name, spd_number, employee_id, position, rank, employee_name, employee_number, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`

	updateAssignee = `
		UPDATE assignees
		SET name = $2, spd_number = $3, employee_id = $4, position = $5, rank = $6, employee_name = $7, employee_number = $8, updated_at = $9
		WHERE id = $1
	`

	findAssigneeByID = `
		SELECT
			a.id, a.business_trip_id, a.name, a.spd_number, a.employee_id, a.position, a.rank, a.employee_name, a.employee_number,
			a.created_at, a.updated_at
		FROM assignees a
		WHERE a.id = $1 AND a.deleted_at IS NULL
	`

	findAssigneesByBusinessTripID = `
		SELECT
			a.id, a.business_trip_id, a.name, a.spd_number, a.employee_id, a.position, a.rank, a.employee_name, a.employee_number,
			a.created_at, a.updated_at
		FROM assignees a
		WHERE a.business_trip_id = $1 AND a.deleted_at IS NULL
		ORDER BY a.created_at
	`

	deleteAssignee = `
		UPDATE assignees
		SET deleted_at = $1
		WHERE id = $2
	`

	insertTransaction = `
		INSERT INTO assignee_transactions (
			id, assignee_id, name, type, subtype, amount, total_night, subtotal,
			description, transport_detail, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`

	updateTransaction = `
		UPDATE assignee_transactions
		SET name = $2, type = $3, subtype = $4, amount = $5, total_night = $6, subtotal = $7,
			description = $8, transport_detail = $9, updated_at = $10
		WHERE id = $1
	`

	findTransactionByID = `
		SELECT
			t.id, t.assignee_id, t.name, t.type, t.subtype, t.amount, t.total_night, t.subtotal,
			t.description, t.transport_detail, t.created_at, t.updated_at
		FROM assignee_transactions t
		WHERE t.id = $1 AND t.deleted_at IS NULL
	`

	findTransactionsByAssigneeID = `
		SELECT
			t.id, t.assignee_id, t.name, t.type, t.subtype, t.amount, t.total_night, t.subtotal,
			t.description, t.transport_detail, t.created_at, t.updated_at
		FROM assignee_transactions t
		WHERE t.assignee_id = $1 AND t.deleted_at IS NULL
		ORDER BY t.created_at
	`

	deleteTransaction = `
		UPDATE assignee_transactions
		SET deleted_at = $1
		WHERE id = $2
	`

	hardDeleteAssigneesByBusinessTripID = `
		DELETE FROM assignees
		WHERE business_trip_id = $1
	`
)

// NewBusinessTripRepository creates a new instance of BusinessTripRepository
func NewBusinessTripRepository(db database.Queryer) repository.BusinessTripRepository {
	// Try to access the underlying *sql.DB if db implements the DB interface
	if dbInterface, ok := db.(database.DB); ok {
		sqlDB := dbInterface.UnderlyingDB()
		return &businessTripRepository{
			db:              db,
			numberGenerator: business_trip_number.NewGenerator(sqlDB),
		}
	}

	// Fallback: create repository without number generator
	// This can happen in tests or when using transactions
	return &businessTripRepository{
		db:              db,
		numberGenerator: nil,
	}
}

type businessTripRepository struct {
	db              database.Queryer
	numberGenerator *business_trip_number.Generator
}

// WithTransaction returns a new repository instance with the given transaction
func (r *businessTripRepository) WithTransaction(tx database.DBTx) repository.BusinessTripRepository {
	return &businessTripRepository{
		db:              tx,
		numberGenerator: r.numberGenerator, // Preserve the number generator from parent
	}
}

// Create creates a new business trip
func (r *businessTripRepository) Create(ctx context.Context, bt *entity.BusinessTrip) (*entity.BusinessTrip, error) {
	if bt.ID == "" {
		bt.ID = uuid.New().String()
	}

	// Generate business trip number
	if r.numberGenerator != nil {
		businessTripNumber, err := r.numberGenerator.GenerateNextNumber(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to generate business trip number: %w", err)
		}
		bt.SetBusinessTripNumber(businessTripNumber)
	} else {
		// Fallback: use UUID-based number if generator is not available
		// Use only 6 characters from UUID to fit within VARCHAR(10) (BT-XXXXXX)
		bt.SetBusinessTripNumber("BT-" + uuid.New().String()[:6])
	}

	var returnedID string
	now := time.Now()

	err := r.db.GetContext(ctx, &returnedID, insertBusinessTrip,
		bt.ID,
		bt.BusinessTripNumber,
		bt.StartDate,
		bt.EndDate,
		bt.ActivityPurpose,
		bt.DestinationCity,
		bt.SPDDate,
		bt.DepartureDate,
		bt.ReturnDate,
		bt.Status,
		bt.DocumentLink,
		now,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create business trip: %w", err)
	}

	if returnedID != bt.ID {
		return nil, fmt.Errorf("returned ID %s does not match expected ID %s", returnedID, bt.ID)
	}
	bt.CreatedAt = now
	bt.UpdatedAt = now

	return bt, nil
}

// GetByID retrieves a business trip by ID with all its related data
func (r *businessTripRepository) GetByID(ctx context.Context, id string) (*entity.BusinessTrip, error) {
	// Get business trip
	var bt entity.BusinessTrip
	err := r.db.GetContext(ctx, &bt, findBusinessTripByID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get business trip: %w", err)
	}

	// Get assignees
	assignees, err := r.GetAssigneesByBusinessTripID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignees: %w", err)
	}

	bt.Assignees = assignees

	return &bt, nil
}

// Update updates a business trip
func (r *businessTripRepository) Update(ctx context.Context, bt *entity.BusinessTrip) (*entity.BusinessTrip, error) {
	now := time.Now()

	res, err := r.db.ExecContext(ctx, updateBusinessTrip,
		bt.ID,
		bt.StartDate,
		bt.EndDate,
		bt.ActivityPurpose,
		bt.DestinationCity,
		bt.SPDDate,
		bt.DepartureDate,
		bt.ReturnDate,
		bt.Status,
		bt.DocumentLink,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update business trip: %w", err)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowAffected == 0 {
		return nil, fmt.Errorf("business trip with ID %s not found", bt.ID)
	}

	bt.UpdatedAt = now

	return bt, nil
}

// Delete soft deletes a business trip
func (r *businessTripRepository) Delete(ctx context.Context, id string) error {
	now := time.Now()

	res, err := r.db.ExecContext(ctx, deleteBusinessTrip, now, id)
	if err != nil {
		return fmt.Errorf("failed to delete business trip: %w", err)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowAffected == 0 {
		return fmt.Errorf("business trip with ID %s not found", id)
	}

	return nil
}

// List retrieves business trips with filtering and pagination using pagination package
func (r *businessTripRepository) List(ctx context.Context, params *pagination.QueryParams) ([]*entity.BusinessTrip, int64, error) {
	// Build count query
	countBuilder := pagination.NewQueryBuilder("SELECT COUNT(*) FROM business_trips")
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
		SELECT
			id, business_trip_number, start_date, end_date, activity_purpose, destination_city,
			spd_date, departure_date, return_date, status, document_link, created_at, updated_at
		FROM business_trips`)

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

	var businessTrips []*entity.BusinessTrip
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var bt entity.BusinessTrip
		if err := rows.StructScan(&bt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan business trip: %w", err)
		}
		businessTrips = append(businessTrips, &bt)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return businessTrips, totalCount, nil
}

// CreateAssignee creates a new assignee
func (r *businessTripRepository) CreateAssignee(ctx context.Context, assignee *entity.Assignee) (*entity.Assignee, error) {
	if assignee.ID == "" {
		assignee.ID = uuid.New().String()
	}

	var returnedID string
	now := time.Now()

	err := r.db.GetContext(ctx, &returnedID, insertAssignee,
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
func (r *businessTripRepository) GetAssigneeByID(ctx context.Context, id string) (*entity.Assignee, error) {
	var assignee entity.Assignee
	err := r.db.GetContext(ctx, &assignee, findAssigneeByID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get assignee: %w", err)
	}

	// Get transactions for this assignee
	transactions, err := r.GetTransactionsByAssigneeID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	assignee.Transactions = transactions

	return &assignee, nil
}

// UpdateAssignee updates an assignee
func (r *businessTripRepository) UpdateAssignee(ctx context.Context, assignee *entity.Assignee) (*entity.Assignee, error) {
	now := time.Now()

	res, err := r.db.ExecContext(ctx, updateAssignee,
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
func (r *businessTripRepository) DeleteAssignee(ctx context.Context, id string) error {
	now := time.Now()

	res, err := r.db.ExecContext(ctx, deleteAssignee, now, id)
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

// GetAssigneesByBusinessTripID retrieves all assignees for a business trip with their transactions
func (r *businessTripRepository) GetAssigneesByBusinessTripID(ctx context.Context, businessTripID string) ([]*entity.Assignee, error) {
	rows, err := r.db.QueryxContext(ctx, findAssigneesByBusinessTripID, businessTripID)
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

		// Get transactions for this assignee
		transactions, err := r.GetTransactionsByAssigneeID(ctx, assignee.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get transactions for assignee %s: %w", assignee.ID, err)
		}

		assignee.Transactions = transactions
		assignees = append(assignees, &assignee)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return assignees, nil
}

// GetAssigneesByBusinessTripIDWithoutTransactions retrieves all assignees for a business trip without loading their transactions
func (r *businessTripRepository) GetAssigneesByBusinessTripIDWithoutTransactions(ctx context.Context, businessTripID string) ([]*entity.Assignee, error) {
	rows, err := r.db.QueryxContext(ctx, findAssigneesByBusinessTripID, businessTripID)
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

// CreateTransaction creates a new transaction
func (r *businessTripRepository) CreateTransaction(ctx context.Context, transaction *entity.Transaction) (*entity.Transaction, error) {
	if transaction.ID == "" {
		transaction.ID = uuid.New().String()
	}

	var returnedID string
	now := time.Now()

	err := r.db.GetContext(ctx, &returnedID, insertTransaction,
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
func (r *businessTripRepository) GetTransactionByID(ctx context.Context, id string) (*entity.Transaction, error) {
	var transaction entity.Transaction
	err := r.db.GetContext(ctx, &transaction, findTransactionByID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &transaction, nil
}

// UpdateTransaction updates a transaction
func (r *businessTripRepository) UpdateTransaction(ctx context.Context, transaction *entity.Transaction) (*entity.Transaction, error) {
	now := time.Now()

	res, err := r.db.ExecContext(ctx, updateTransaction,
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
func (r *businessTripRepository) DeleteTransaction(ctx context.Context, id string) error {
	now := time.Now()

	res, err := r.db.ExecContext(ctx, deleteTransaction, now, id)
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
func (r *businessTripRepository) GetTransactionsByAssigneeID(ctx context.Context, assigneeID string) ([]*entity.Transaction, error) {
	rows, err := r.db.QueryxContext(ctx, findTransactionsByAssigneeID, assigneeID)
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

// DeleteAssigneesByBusinessTripID deletes all assignees for a business trip
func (r *businessTripRepository) DeleteAssigneesByBusinessTripID(ctx context.Context, businessTripID string) error {
	_, err := r.db.ExecContext(ctx, hardDeleteAssigneesByBusinessTripID, businessTripID)
	if err != nil {
		return fmt.Errorf("failed to delete assignees by business trip ID: %w", err)
	}
	return nil
}

// DeleteTransactionsByAssigneeIDs deletes transactions by multiple assignee IDs
func (r *businessTripRepository) DeleteTransactionsByAssigneeIDs(ctx context.Context, assigneeIDs []string) error {
	if len(assigneeIDs) == 0 {
		return nil
	}

	// Build the query dynamically for each assignee ID
	placeholders := make([]string, len(assigneeIDs))
	args := make([]interface{}, len(assigneeIDs))

	for i, id := range assigneeIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf("DELETE FROM assignee_transactions WHERE assignee_id IN (%s)", strings.Join(placeholders, ","))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete transactions by assignee IDs: %w", err)
	}
	return nil
}

// GetDestinationStats gets destination statistics for the dashboard
func (r *businessTripRepository) GetDestinationStats(ctx context.Context, startDate, endDate *time.Time, destination string) ([]*repository.DestinationData, error) {
	// For now, return empty slice as a placeholder implementation
	// TODO: Implement actual destination statistics query
	return []*repository.DestinationData{}, nil
}

// GetMonthlyStats gets monthly statistics for the dashboard
func (r *businessTripRepository) GetMonthlyStats(ctx context.Context, startDate, endDate time.Time, destination string) ([]*repository.MonthlyData, error) {
	// For now, return empty slice as a placeholder implementation
	// TODO: Implement actual monthly statistics query
	return []*repository.MonthlyData{}, nil
}

// GetRecentWithSummary gets recent business trips with summary data
func (r *businessTripRepository) GetRecentWithSummary(ctx context.Context, limit int) ([]*repository.RecentBusinessTripData, error) {
	// For now, return empty slice as a placeholder implementation
	// TODO: Implement actual recent business trips with summary query
	return []*repository.RecentBusinessTripData{}, nil
}

// GetStatusCounts gets status counts for the dashboard
func (r *businessTripRepository) GetStatusCounts(ctx context.Context, startDate, endDate *time.Time, destination string) (*repository.StatusCounts, error) {
	// For now, return empty counts as a placeholder implementation
	// TODO: Implement actual status counts query
	return &repository.StatusCounts{
		Total:     0,
		Draft:     0,
		Ongoing:   0,
		Completed: 0,
		Canceled:  0,
	}, nil
}

// GetTotalCost gets total cost for the dashboard
func (r *businessTripRepository) GetTotalCost(ctx context.Context, startDate, endDate *time.Time, destination string) (float64, error) {
	// For now, return 0 as a placeholder implementation
	// TODO: Implement actual total cost query
	return 0.0, nil
}

// GetTotalCount gets total count for the dashboard
func (r *businessTripRepository) GetTotalCount(ctx context.Context, startDate, endDate *time.Time) (int64, error) {
	// For now, return 0 as a placeholder implementation
	// TODO: Implement actual total count query
	return 0, nil
}

// GetTypeStats gets transaction type statistics for the dashboard
func (r *businessTripRepository) GetTypeStats(ctx context.Context, startDate, endDate *time.Time) ([]*repository.TransactionTypeData, error) {
	// For now, return empty slice as a placeholder implementation
	// TODO: Implement actual type statistics query
	return []*repository.TransactionTypeData{}, nil
}

// GetUpcomingCount gets upcoming business trips count for the dashboard
func (r *businessTripRepository) GetUpcomingCount(ctx context.Context) (int64, error) {
	// For now, return 0 as a placeholder implementation
	// TODO: Implement actual upcoming count query
	return 0, nil
}
