package repository

import (
	"context"

	"sandbox/internal/domain/entity"
)

// BusinessTripHistoryRepository defines the interface for business trip history repository
type BusinessTripHistoryRepository interface {
	// Create creates a new business trip history record
	Create(ctx context.Context, history *entity.BusinessTripHistory) error

	// FindByBusinessTripID retrieves all history records for a specific business trip
	FindByBusinessTripID(ctx context.Context, businessTripID string) ([]*entity.BusinessTripHistory, error)

	// FindByID retrieves a history record by its ID
	FindByID(ctx context.Context, id string) (*entity.BusinessTripHistory, error)

	// FindByBusinessTripIDAndChangeType retrieves history records by business trip ID and change type
	FindByBusinessTripIDAndChangeType(ctx context.Context, businessTripID string, changeType entity.BusinessTripHistoryChangeType) ([]*entity.BusinessTripHistory, error)
}
