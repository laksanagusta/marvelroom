package repository

import (
	"context"

	"sandbox/internal/domain/entity"
	"sandbox/pkg/pagination"
)

// BusinessTripRepository defines the interface for business trip data operations
type BusinessTripRepository interface {
	// Business Trip operations
	Create(ctx context.Context, bt *entity.BusinessTrip) (*entity.BusinessTrip, error)
	GetByID(ctx context.Context, id string) (*entity.BusinessTrip, error)
	Update(ctx context.Context, bt *entity.BusinessTrip) (*entity.BusinessTrip, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, params *pagination.QueryParams) ([]*entity.BusinessTrip, int64, error)

	// Assignee operations
	CreateAssignee(ctx context.Context, assignee *entity.Assignee) (*entity.Assignee, error)
	GetAssigneeByID(ctx context.Context, id string) (*entity.Assignee, error)
	UpdateAssignee(ctx context.Context, assignee *entity.Assignee) (*entity.Assignee, error)
	DeleteAssignee(ctx context.Context, id string) error
	GetAssigneesByBusinessTripID(ctx context.Context, businessTripID string) ([]*entity.Assignee, error)
	GetAssigneesByBusinessTripIDWithoutTransactions(ctx context.Context, businessTripID string) ([]*entity.Assignee, error)
	DeleteAssigneesByBusinessTripID(ctx context.Context, businessTripID string) error

	// Transaction operations
	CreateTransaction(ctx context.Context, transaction *entity.Transaction) (*entity.Transaction, error)
	GetTransactionByID(ctx context.Context, id string) (*entity.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *entity.Transaction) (*entity.Transaction, error)
	DeleteTransaction(ctx context.Context, id string) error
	GetTransactionsByAssigneeID(ctx context.Context, assigneeID string) ([]*entity.Transaction, error)
	DeleteTransactionsByAssigneeIDs(ctx context.Context, assigneeIDs []string) error
}

