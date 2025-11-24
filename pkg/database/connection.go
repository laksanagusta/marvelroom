package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewConnectionx creates a new sqlx database connection with proper configuration
func NewConnectionx(dsn string) (*sqlx.DB, error) {
	// Use sqlx.Connect which opens and pings the database in one step.
	dbx, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sqlx database: %w", err)
	}

	// Apply the connection pool settings directly to the sqlx.DB object.
	// The underlying *sql.DB is accessible.
	dbx.SetMaxOpenConns(25)
	dbx.SetMaxIdleConns(5)
	dbx.SetConnMaxLifetime(5 * time.Minute)
	dbx.SetConnMaxIdleTime(5 * time.Minute)

	log.Printf("✅ Database (sqlx) connected successfully")
	return dbx, nil
}

// TestConnection tests if the database is reachable
func TestConnection(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database for testing: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}
