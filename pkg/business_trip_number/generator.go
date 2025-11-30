package business_trip_number

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// Generator handles business trip number generation
type Generator struct {
	db *sql.DB
}

// NewGenerator creates a new business trip number generator
func NewGenerator(db *sql.DB) *Generator {
	return &Generator{db: db}
}

// SetDatabase allows updating the database connection (useful for testing)
func (g *Generator) SetDatabase(db *sql.DB) {
	g.db = db
}

// GenerateNextNumber generates the next business trip number in format BT-XXXXXX
func (g *Generator) GenerateNextNumber(ctx context.Context) (string, error) {
	// Start a transaction for atomic operation
	tx, err := g.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get the current maximum sequence number
	var maxSeq int
	query := `
		SELECT COALESCE(MAX(CAST(SUBSTRING(business_trip_number FROM 4) AS INTEGER)), 0)
		FROM business_trips
		WHERE business_trip_number LIKE 'BT-%'
		AND business_trip_number ~ '^BT-[0-9]{6}$'
		AND deleted_at IS NULL
	`

	err = tx.QueryRowContext(ctx, query).Scan(&maxSeq)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to get max sequence number: %w", err)
	}

	// Increment the sequence
	nextSeq := maxSeq + 1
	nextNumber := fmt.Sprintf("BT-%06d", nextSeq)

	// Verify this number doesn't exist (double-check)
	var exists bool
	checkQuery := `
		SELECT EXISTS(
			SELECT 1 FROM business_trips
			WHERE business_trip_number = $1
			AND deleted_at IS NULL
		)
	`
	err = tx.QueryRowContext(ctx, checkQuery, nextNumber).Scan(&exists)
	if err != nil {
		return "", fmt.Errorf("failed to check number existence: %w", err)
	}

	if exists {
		// If by some chance it exists, try again recursively
		log.Printf("Business trip number %s already exists, trying next", nextNumber)
		return g.GenerateNextNumber(ctx)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nextNumber, nil
}
