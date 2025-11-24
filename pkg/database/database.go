package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Queryer defines the interface for database operations
type Queryer interface {
	// Query methods
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// Exec methods
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// Transaction methods
	Rebind(query string) string
}

// DBTx defines the interface for database transactions
type DBTx interface {
	Queryer
	Commit() error
	Rollback() error
}

// DB extends Queryer with additional database-specific methods
type DB interface {
	Queryer
	BeginTx(ctx context.Context, opts *sql.TxOptions) (DBTx, error)
	WithTransaction(ctx context.Context, fn func(ctx context.Context, tx DBTx) error) error
	Close() error
	UnderlyingDB() *sql.DB
}

// NewDB creates a new DB instance from sqlx.DB
func NewDB(db *sqlx.DB) DB {
	return &dbWrapper{db: db}
}

// NewTx creates a new DBTx instance from sqlx.Tx
func NewTx(tx *sqlx.Tx) DBTx {
	return &txWrapper{tx: tx}
}

// dbWrapper implements DB interface
type dbWrapper struct {
	db *sqlx.DB
}

func (w *dbWrapper) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return w.db.QueryxContext(ctx, query, args...)
}

func (w *dbWrapper) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return w.db.QueryRowxContext(ctx, query, args...)
}

func (w *dbWrapper) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return w.db.GetContext(ctx, dest, query, args...)
}

func (w *dbWrapper) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return w.db.SelectContext(ctx, dest, query, args...)
}

func (w *dbWrapper) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return w.db.ExecContext(ctx, query, args...)
}

func (w *dbWrapper) Rebind(query string) string {
	return w.db.Rebind(query)
}

func (w *dbWrapper) BeginTx(ctx context.Context, opts *sql.TxOptions) (DBTx, error) {
	tx, err := w.db.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return NewTx(tx), nil
}

func (w *dbWrapper) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx DBTx) error) error {
	tx, err := w.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(ctx, tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback failed: %w", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func (w *dbWrapper) UnderlyingDB() *sql.DB {
	return w.db.DB
}

func (w *dbWrapper) Close() error {
	return w.db.Close()
}

// txWrapper implements DBTx interface
type txWrapper struct {
	tx *sqlx.Tx
}

func (w *txWrapper) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return w.tx.QueryxContext(ctx, query, args...)
}

func (w *txWrapper) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return w.tx.QueryRowxContext(ctx, query, args...)
}

func (w *txWrapper) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return w.tx.GetContext(ctx, dest, query, args...)
}

func (w *txWrapper) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return w.tx.SelectContext(ctx, dest, query, args...)
}

func (w *txWrapper) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return w.tx.ExecContext(ctx, query, args...)
}

func (w *txWrapper) Rebind(query string) string {
	return w.tx.Rebind(query)
}

func (w *txWrapper) Commit() error {
	return w.tx.Commit()
}

func (w *txWrapper) Rollback() error {
	return w.tx.Rollback()
}