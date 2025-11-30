package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"sandbox/internal/domain/entity"
	"sandbox/pkg/pagination"
)

// StatusCounts represents count of business trips by status
type StatusCounts struct {
	Total     int64 `json:"total"`
	Draft     int64 `json:"draft"`
	Ongoing   int64 `json:"ongoing"`
	Completed  int64 `json:"completed"`
	Canceled   int64 `json:"canceled"`
}

// MonthlyData represents monthly statistics data
type MonthlyData struct {
	Month           string  `json:"month"`
	Year            int     `json:"year"`
	TotalTrips      int64   `json:"total_trips"`
	CompletedTrips   int64   `json:"completed_trips"`
	TotalCost        float64 `json:"total_cost"`
	TopDestination   string  `json:"top_destination"`
}

// DestinationData represents destination statistics data
type DestinationData struct {
	Destination   string    `json:"destination"`
	TotalTrips    int64     `json:"total_trips"`
	CompletedTrips int64     `json:"completed_trips"`
	TotalCost      float64    `json:"total_cost"`
	LastTripDate   time.Time  `json:"last_trip_date"`
}

// RecentBusinessTripData represents recent business trip with summary
type RecentBusinessTripData struct {
	ID            uuid.UUID `json:"id"`
	BusinessTripNumber string     `json:"business_trip_number"`
	ActivityPurpose  string     `json:"activity_purpose"`
	DestinationCity  string     `json:"destination_city"`
	StartDate       time.Time   `json:"start_date"`
	EndDate         time.Time   `json:"end_date"`
	Status          entity.BusinessTripStatus `json:"status"`
	AssigneeCount   int64       `json:"assignee_count"`
	TotalCost        float64     `json:"total_cost"`
}

// TransactionTypeData represents transaction type statistics for dashboard
type TransactionTypeData struct {
	TransactionType   string  `json:"transaction_type"`
	TotalTransactions  int64   `json:"total_transactions"`
	TotalAmount       float64 `json:"total_amount"`
	AverageAmount     float64 `json:"average_amount"`
}

// BusinessTripRepository defines the interface for business trip data operations
type BusinessTripRepository interface {
	// Business Trip operations
	Create(ctx context.Context, bt *entity.BusinessTrip) (*entity.BusinessTrip, error)
	GetByID(ctx context.Context, id string) (*entity.BusinessTrip, error)
	Update(ctx context.Context, bt *entity.BusinessTrip) (*entity.BusinessTrip, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, params *pagination.QueryParams) ([]*entity.BusinessTrip, int64, error)

	// Dashboard operations
	GetStatusCounts(ctx context.Context, startDate, endDate *time.Time, destination string) (*StatusCounts, error)
	GetTotalCost(ctx context.Context, startDate, endDate *time.Time, destination string) (float64, error)
	GetMonthlyStats(ctx context.Context, startDate, endDate time.Time, destination string) ([]*MonthlyData, error)
	GetDestinationStats(ctx context.Context, startDate, endDate *time.Time, destination string) ([]*DestinationData, error)
	GetUpcomingCount(ctx context.Context) (int64, error)
	GetRecentWithSummary(ctx context.Context, limit int) ([]*RecentBusinessTripData, error)

	// Assignee operations (for dashboard)
	GetTotalCount(ctx context.Context, startDate, endDate *time.Time) (int64, error)

	// Transaction operations (for dashboard)
	GetTypeStats(ctx context.Context, startDate, endDate *time.Time) ([]*TransactionTypeData, error)

	// Transaction operations
	CreateTransaction(ctx context.Context, transaction *entity.Transaction) (*entity.Transaction, error)
	GetTransactionByID(ctx context.Context, id string) (*entity.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *entity.Transaction) (*entity.Transaction, error)
	DeleteTransaction(ctx context.Context, id string) error
	GetTransactionsByAssigneeID(ctx context.Context, assigneeID string) ([]*entity.Transaction, error)
	DeleteTransactionsByAssigneeIDs(ctx context.Context, assigneeIDs []string) error
}
