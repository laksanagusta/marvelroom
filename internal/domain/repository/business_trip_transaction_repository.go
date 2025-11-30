package repository

import (
	"context"
	"time"

	"sandbox/internal/domain/entity"
)

// BusinessTripTransactionRepository defines the interface for business trip transaction data operations
type BusinessTripTransactionRepository interface {
	// Transaction operations
	CreateTransaction(ctx context.Context, transaction *entity.Transaction) (*entity.Transaction, error)
	GetTransactionByID(ctx context.Context, id string) (*entity.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *entity.Transaction) (*entity.Transaction, error)
	DeleteTransaction(ctx context.Context, id string) error
	GetTransactionsByAssigneeID(ctx context.Context, assigneeID string) ([]*entity.Transaction, error)
	DeleteTransactionsByAssigneeIDs(ctx context.Context, assigneeIDs []string) error
	GetTotalCount(ctx context.Context, startDate, endDate *time.Time) (int64, error)

	// Dashboard operations
	GetTypeStats(ctx context.Context, startDate, endDate *time.Time) ([]*TransactionTypeData, error)
}