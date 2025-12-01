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

	// Verificator SQL queries
	insertVerificator = `
		INSERT INTO business_trip_verificators (
			id, business_trip_id, user_id, user_name, employee_number, position, status,
			verification_notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	findVerificatorByID = `
		SELECT
			v.id, v.business_trip_id, v.user_id, v.user_name, v.employee_number, v.position,
			v.status, v.verified_at, v.verification_notes, v.created_at, v.updated_at
		FROM business_trip_verificators v
		WHERE v.id = $1 AND v.deleted_at IS NULL
	`

	findVerificatorsByBusinessTripID = `
		SELECT
			v.id, v.business_trip_id, v.user_id, v.user_name, v.employee_number, v.position,
			v.status, v.verified_at, v.verification_notes, v.created_at, v.updated_at
		FROM business_trip_verificators v
		WHERE v.business_trip_id = $1 AND v.deleted_at IS NULL
		ORDER BY v.created_at
	`

	findVerificators = `
		SELECT
			v.id, v.business_trip_id, v.user_id, v.user_name, v.employee_number, v.position,
			v.status, v.verified_at, v.verification_notes, v.created_at, v.updated_at,
			bt.business_trip_number, bt.start_date, bt.end_date, bt.activity_purpose,
			bt.destination_city, bt.spd_date, bt.departure_date, bt.return_date,
			bt.status as business_trip_status, bt.document_link
		FROM business_trip_verificators v
		LEFT JOIN business_trips bt ON v.business_trip_id = bt.id
	`

	findVerificatorByBusinessTripIDAndUserID = `
		SELECT
			v.id, v.business_trip_id, v.user_id, v.user_name, v.employee_number, v.position,
			v.status, v.verified_at, v.verification_notes, v.created_at, v.updated_at
		FROM business_trip_verificators v
		WHERE v.business_trip_id = $1 AND v.user_id = $2 AND v.deleted_at IS NULL
		ORDER BY v.created_at
		LIMIT 1
	`

	updateVerificator = `
		UPDATE business_trip_verificators
		SET status = $2, verification_notes = $3, verified_at = $4, updated_at = $5
		WHERE id = $1
	`

	deleteVerificator = `
		UPDATE business_trip_verificators
		SET deleted_at = $1
		WHERE id = $2
	`

	deleteVerificatorsByBusinessTripID = `
		UPDATE business_trip_verificators
		SET deleted_at = $1
		WHERE business_trip_id = $2
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

	// Get verificators
	verificators, err := r.GetVerificatorsByBusinessTripID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get verificators: %w", err)
	}

	bt.Verificators = verificators

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

	// First, collect all assignees without transactions
	var assignees []*entity.Assignee
	for rows.Next() {
		var assignee entity.Assignee
		err := rows.StructScan(&assignee)
		if err != nil {
			return nil, fmt.Errorf("failed to scan assignee: %w", err)
		}
		assignees = append(assignees, &assignee)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Close rows before executing nested queries
	rows.Close()

	// Now fetch transactions for each assignee
	for _, assignee := range assignees {
		transactions, err := r.GetTransactionsByAssigneeID(ctx, assignee.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get transactions for assignee %s: %w", assignee.ID, err)
		}
		assignee.Transactions = transactions
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
	query := `
		SELECT
			bt.destination_city,
			COUNT(*) as total_trips,
			COUNT(*) FILTER (WHERE bt.status = 'completed') as completed_trips,
			COALESCE(SUM(t.subtotal), 0) as total_cost,
			MAX(bt.start_date) as last_trip_date
		FROM business_trips bt
		LEFT JOIN assignees a ON bt.id = a.business_trip_id AND a.deleted_at IS NULL
		LEFT JOIN assignee_transactions t ON a.id = t.assignee_id AND t.deleted_at IS NULL
		WHERE bt.deleted_at IS NULL
	`

	args := []interface{}{}
	argIndex := 1

	if startDate != nil {
		query += fmt.Sprintf(" AND bt.start_date >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}

	if endDate != nil {
		query += fmt.Sprintf(" AND bt.end_date <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}

	if destination != "" {
		query += fmt.Sprintf(" AND bt.destination_city ILIKE $%d", argIndex)
		args = append(args, "%"+destination+"%")
		argIndex++
	}

	query += " GROUP BY bt.destination_city ORDER BY total_cost DESC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get destination stats: %w", err)
	}
	defer rows.Close()

	var stats []*repository.DestinationData
	for rows.Next() {
		var stat repository.DestinationData
		err := rows.Scan(
			&stat.Destination,
			&stat.TotalTrips,
			&stat.CompletedTrips,
			&stat.TotalCost,
			&stat.LastTripDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan destination stat: %w", err)
		}
		stats = append(stats, &stat)
	}

	return stats, nil
}

// GetMonthlyStats gets monthly statistics for the dashboard
func (r *businessTripRepository) GetMonthlyStats(ctx context.Context, startDate, endDate time.Time, destination string) ([]*repository.MonthlyData, error) {
	query := `
		SELECT
			TO_CHAR(bt.start_date, 'Month') as month,
			EXTRACT(YEAR FROM bt.start_date) as year,
			COUNT(*) as total_trips,
			COUNT(*) FILTER (WHERE bt.status = 'completed') as completed_trips,
			COALESCE(SUM(t.subtotal), 0) as total_cost,
			MODE() WITHIN GROUP (ORDER BY bt.destination_city) as top_destination
		FROM business_trips bt
		LEFT JOIN assignees a ON bt.id = a.business_trip_id AND a.deleted_at IS NULL
		LEFT JOIN assignee_transactions t ON a.id = t.assignee_id AND t.deleted_at IS NULL
		WHERE bt.deleted_at IS NULL
		AND bt.start_date >= $1 AND bt.start_date <= $2
	`

	args := []interface{}{startDate, endDate}
	argIndex := 3

	if destination != "" {
		query += fmt.Sprintf(" AND bt.destination_city ILIKE $%d", argIndex)
		args = append(args, "%"+destination+"%")
		argIndex++
	}

	query += " GROUP BY TO_CHAR(bt.start_date, 'Month'), EXTRACT(YEAR FROM bt.start_date) ORDER BY year DESC, month"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly stats: %w", err)
	}
	defer rows.Close()

	var stats []*repository.MonthlyData
	for rows.Next() {
		var stat repository.MonthlyData
		err := rows.Scan(
			&stat.Month,
			&stat.Year,
			&stat.TotalTrips,
			&stat.CompletedTrips,
			&stat.TotalCost,
			&stat.TopDestination,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan monthly stat: %w", err)
		}
		stats = append(stats, &stat)
	}

	return stats, nil
}

// GetRecentWithSummary gets recent business trips with summary data
func (r *businessTripRepository) GetRecentWithSummary(ctx context.Context, limit int) ([]*repository.RecentBusinessTripData, error) {
	query := `
		SELECT
			bt.id,
			bt.business_trip_number,
			bt.activity_purpose,
			bt.destination_city,
			bt.start_date,
			bt.end_date,
			bt.status,
			COUNT(a.id) as assignee_count,
			COALESCE(SUM(t.subtotal), 0) as total_cost
		FROM business_trips bt
		LEFT JOIN assignees a ON bt.id = a.business_trip_id AND a.deleted_at IS NULL
		LEFT JOIN assignee_transactions t ON a.id = t.assignee_id AND t.deleted_at IS NULL
		WHERE bt.deleted_at IS NULL
		GROUP BY bt.id, bt.business_trip_number, bt.activity_purpose, bt.destination_city,
			bt.start_date, bt.end_date, bt.status
		ORDER BY bt.created_at DESC
		LIMIT $1
	`

	rows, err := r.db.QueryxContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent business trips with summary: %w", err)
	}
	defer rows.Close()

	var trips []*repository.RecentBusinessTripData
	for rows.Next() {
		var trip repository.RecentBusinessTripData
		err := rows.Scan(
			&trip.ID,
			&trip.BusinessTripNumber,
			&trip.ActivityPurpose,
			&trip.DestinationCity,
			&trip.StartDate,
			&trip.EndDate,
			&trip.Status,
			&trip.AssigneeCount,
			&trip.TotalCost,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recent business trip: %w", err)
		}
		trips = append(trips, &trip)
	}

	return trips, nil
}

// GetStatusCounts gets status counts for the dashboard
func (r *businessTripRepository) GetStatusCounts(ctx context.Context, startDate, endDate *time.Time, destination string) (*repository.StatusCounts, error) {
	query := `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'draft') as draft,
			COUNT(*) FILTER (WHERE status = 'ongoing') as ongoing,
			COUNT(*) FILTER (WHERE status = 'completed') as completed,
			COUNT(*) FILTER (WHERE status = 'canceled') as canceled
		FROM business_trips
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	argIndex := 1

	if startDate != nil {
		query += fmt.Sprintf(" AND start_date >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}

	if endDate != nil {
		query += fmt.Sprintf(" AND end_date <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}

	if destination != "" {
		query += fmt.Sprintf(" AND destination_city ILIKE $%d", argIndex)
		args = append(args, "%"+destination+"%")
		argIndex++
	}

	var counts repository.StatusCounts
	err := r.db.QueryRowxContext(ctx, query, args...).Scan(
		&counts.Total,
		&counts.Draft,
		&counts.Ongoing,
		&counts.Completed,
		&counts.Canceled,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get status counts: %w", err)
	}

	return &counts, nil
}

// GetTotalCost gets total cost for the dashboard
func (r *businessTripRepository) GetTotalCost(ctx context.Context, startDate, endDate *time.Time, destination string) (float64, error) {
	query := `
		SELECT COALESCE(SUM(t.subtotal), 0) as total_cost
		FROM business_trips bt
		LEFT JOIN assignees a ON bt.id = a.business_trip_id AND a.deleted_at IS NULL
		LEFT JOIN assignee_transactions t ON a.id = t.assignee_id AND t.deleted_at IS NULL
		WHERE bt.deleted_at IS NULL
	`

	args := []interface{}{}
	argIndex := 1

	if startDate != nil {
		query += fmt.Sprintf(" AND bt.start_date >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}

	if endDate != nil {
		query += fmt.Sprintf(" AND bt.end_date <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}

	if destination != "" {
		query += fmt.Sprintf(" AND bt.destination_city ILIKE $%d", argIndex)
		args = append(args, "%"+destination+"%")
		argIndex++
	}

	var totalCost float64
	err := r.db.QueryRowxContext(ctx, query, args...).Scan(&totalCost)
	if err != nil {
		return 0, fmt.Errorf("failed to get total cost: %w", err)
	}

	return totalCost, nil
}

// GetTotalCount gets total count for the dashboard
func (r *businessTripRepository) GetTotalCount(ctx context.Context, startDate, endDate *time.Time) (int64, error) {
	query := `
		SELECT COUNT(DISTINCT a.id) as total_count
		FROM business_trips bt
		LEFT JOIN assignees a ON bt.id = a.business_trip_id AND a.deleted_at IS NULL
		WHERE bt.deleted_at IS NULL
	`

	args := []interface{}{}
	argIndex := 1

	if startDate != nil {
		query += fmt.Sprintf(" AND bt.start_date >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}

	if endDate != nil {
		query += fmt.Sprintf(" AND bt.end_date <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}

	var totalCount int64
	err := r.db.QueryRowxContext(ctx, query, args...).Scan(&totalCount)
	if err != nil {
		return 0, fmt.Errorf("failed to get total count: %w", err)
	}

	return totalCount, nil
}

// GetTypeStats gets transaction type statistics for the dashboard
func (r *businessTripRepository) GetTypeStats(ctx context.Context, startDate, endDate *time.Time) ([]*repository.TransactionTypeData, error) {
	query := `
		SELECT
			t.type,
			COUNT(*) as total_transactions,
			COALESCE(SUM(t.subtotal), 0) as total_amount,
			COALESCE(AVG(t.subtotal), 0) as average_amount
		FROM business_trips bt
		LEFT JOIN assignees a ON bt.id = a.business_trip_id AND a.deleted_at IS NULL
		LEFT JOIN assignee_transactions t ON a.id = t.assignee_id AND t.deleted_at IS NULL
		WHERE bt.deleted_at IS NULL
		AND t.type IS NOT NULL
	`

	args := []interface{}{}
	argIndex := 1

	if startDate != nil {
		query += fmt.Sprintf(" AND bt.start_date >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}

	if endDate != nil {
		query += fmt.Sprintf(" AND bt.end_date <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}

	query += " GROUP BY t.type ORDER BY total_amount DESC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction type stats: %w", err)
	}
	defer rows.Close()

	var stats []*repository.TransactionTypeData
	for rows.Next() {
		var stat repository.TransactionTypeData
		err := rows.Scan(
			&stat.TransactionType,
			&stat.TotalTransactions,
			&stat.TotalAmount,
			&stat.AverageAmount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction type stat: %w", err)
		}
		stats = append(stats, &stat)
	}

	return stats, nil
}

// GetUpcomingCount gets upcoming business trips count for the dashboard
func (r *businessTripRepository) GetUpcomingCount(ctx context.Context) (int64, error) {
	query := `
		SELECT COUNT(*) as upcoming_count
		FROM business_trips
		WHERE deleted_at IS NULL
		AND status IN ('draft', 'ongoing')
		AND start_date > NOW()
	`

	var count int64
	err := r.db.QueryRowxContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get upcoming count: %w", err)
	}

	return count, nil
}

// CreateVerificator creates a new verificator
func (r *businessTripRepository) CreateVerificator(ctx context.Context, verificator *entity.Verificator) (*entity.Verificator, error) {
	if verificator.ID == "" {
		verificator.ID = uuid.New().String()
	}

	var returnedID string
	now := time.Now()

	err := r.db.GetContext(ctx, &returnedID, insertVerificator,
		verificator.ID,
		verificator.BusinessTripID,
		verificator.UserID,
		verificator.UserName,
		verificator.EmployeeNumber,
		verificator.Position,
		verificator.Status,
		verificator.VerificationNotes,
		now,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create verificator: %w", err)
	}

	if returnedID != verificator.ID {
		return nil, fmt.Errorf("returned ID %s does not match expected ID %s", returnedID, verificator.ID)
	}

	verificator.CreatedAt = now
	verificator.UpdatedAt = now

	return verificator, nil
}

// GetVerificatorByID retrieves a verificator by ID
func (r *businessTripRepository) GetVerificatorByID(ctx context.Context, id string) (*entity.Verificator, error) {
	var verificator entity.Verificator
	err := r.db.GetContext(ctx, &verificator, findVerificatorByID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get verificator: %w", err)
	}

	return &verificator, nil
}

// GetVerificatorsByBusinessTripID retrieves all verificators for a business trip
func (r *businessTripRepository) GetVerificatorsByBusinessTripID(ctx context.Context, businessTripID string) ([]*entity.Verificator, error) {
	rows, err := r.db.QueryxContext(ctx, findVerificatorsByBusinessTripID, businessTripID)
	if err != nil {
		return nil, fmt.Errorf("failed to query verificators: %w", err)
	}
	defer rows.Close()

	var verificators []*entity.Verificator
	for rows.Next() {
		var verificator entity.Verificator
		err := rows.StructScan(&verificator)
		if err != nil {
			return nil, fmt.Errorf("failed to scan verificator: %w", err)
		}
		verificators = append(verificators, &verificator)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return verificators, nil
}

// GetVerificatorByBusinessTripIDAndUserID retrieves a verificator by business trip ID and user ID
func (r *businessTripRepository) GetVerificatorByBusinessTripIDAndUserID(ctx context.Context, businessTripID, userID string) (*entity.Verificator, error) {
	var verificator entity.Verificator
	err := r.db.GetContext(ctx, &verificator, findVerificatorByBusinessTripIDAndUserID, businessTripID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get verificator: %w", err)
	}

	return &verificator, nil
}

// UpdateVerificator updates a verificator
func (r *businessTripRepository) UpdateVerificator(ctx context.Context, verificator *entity.Verificator) (*entity.Verificator, error) {
	now := time.Now()

	res, err := r.db.ExecContext(ctx, updateVerificator,
		verificator.ID,
		verificator.Status,
		verificator.VerificationNotes,
		verificator.VerifiedAt,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update verificator: %w", err)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowAffected == 0 {
		return nil, fmt.Errorf("verificator with ID %s not found", verificator.ID)
	}

	verificator.UpdatedAt = now

	return verificator, nil
}

// DeleteVerificator soft deletes a verificator
func (r *businessTripRepository) DeleteVerificator(ctx context.Context, id string) error {
	now := time.Now()

	res, err := r.db.ExecContext(ctx, deleteVerificator, now, id)
	if err != nil {
		return fmt.Errorf("failed to delete verificator: %w", err)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowAffected == 0 {
		return fmt.Errorf("verificator with ID %s not found", id)
	}

	return nil
}

// DeleteVerificatorsByBusinessTripID soft deletes all verificators for a business trip
func (r *businessTripRepository) DeleteVerificatorsByBusinessTripID(ctx context.Context, businessTripID string) error {
	now := time.Now()

	_, err := r.db.ExecContext(ctx, deleteVerificatorsByBusinessTripID, now, businessTripID)
	if err != nil {
		return fmt.Errorf("failed to delete verificators by business trip ID: %w", err)
	}

	return nil
}

// ListVerificators retrieves verificators with filtering and pagination using pagination package
func (r *businessTripRepository) ListVerificators(ctx context.Context, params *pagination.QueryParams) ([]*entity.VerificatorWithBusinessTrip, int64, error) {
	// Build count query
	countBuilder := pagination.NewQueryBuilder("SELECT COUNT(*) FROM business_trip_verificators v LEFT JOIN business_trips bt ON v.business_trip_id = bt.id")
	for _, filter := range params.Filters {
		if err := countBuilder.AddFilter(filter); err != nil {
			return nil, 0, err
		}
	}

	// Always include deleted_at filter
	countBuilder.AddFilter(pagination.Filter{
		Field:    "v.deleted_at",
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
	queryBuilder := pagination.NewQueryBuilder(findVerificators)

	for _, filter := range params.Filters {
		if err := queryBuilder.AddFilter(filter); err != nil {
			return nil, 0, err
		}
	}

	// Always include deleted_at filter
	queryBuilder.AddFilter(pagination.Filter{
		Field:    "v.deleted_at",
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

	var verificators []*entity.VerificatorWithBusinessTrip
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var verificator entity.VerificatorWithBusinessTrip
		if err := rows.StructScan(&verificator); err != nil {
			return nil, 0, fmt.Errorf("failed to scan verificator: %w", err)
		}
		verificators = append(verificators, &verificator)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return verificators, totalCount, nil
}
