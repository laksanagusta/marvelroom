package repository

import (
	"context"
	"time"

	"sandbox/internal/domain/entity"
)

// AssigneeRepository defines the interface for assignee data operations
type AssigneeRepository interface {
	// Assignee operations
	Create(ctx context.Context, assignee *entity.Assignee) (*entity.Assignee, error)
	GetAssigneeByID(ctx context.Context, id string) (*entity.Assignee, error)
	UpdateAssignee(ctx context.Context, assignee *entity.Assignee) (*entity.Assignee, error)
	DeleteAssignee(ctx context.Context, id string) error
	GetAssigneesByBusinessTripID(ctx context.Context, businessTripID string) ([]*entity.Assignee, error)
	GetAssigneesByBusinessTripIDWithoutTransactions(ctx context.Context, businessTripID string) ([]*entity.Assignee, error)
	DeleteAssigneesByBusinessTripID(ctx context.Context, businessTripID string) error

	// Dashboard operations
	GetTotalCount(ctx context.Context, startDate, endDate *time.Time) (int64, error)
}