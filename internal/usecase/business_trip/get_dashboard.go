package business_trip

import (
	"context"
	"time"

	"sandbox/internal/domain/repository"
)

// GetDashboardUseCase handles retrieving dashboard data for business trips
type GetDashboardUseCase struct {
	businessTripRepo repository.BusinessTripRepository
	assigneeRepo     repository.AssigneeRepository
	transactionRepo  repository.BusinessTripTransactionRepository
}

// NewGetDashboardUseCase creates a new use case instance
func NewGetDashboardUseCase(
	businessTripRepo repository.BusinessTripRepository,
	assigneeRepo repository.AssigneeRepository,
	transactionRepo repository.BusinessTripTransactionRepository,
) *GetDashboardUseCase {
	return &GetDashboardUseCase{
		businessTripRepo: businessTripRepo,
		assigneeRepo:     assigneeRepo,
		transactionRepo:  transactionRepo,
	}
}

// GetDashboardRequest represents the request parameters for dashboard data
type GetDashboardRequest struct {
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Destination string     `json:"destination,omitempty"`
	Status      string     `json:"status,omitempty"`
	Limit       int        `json:"limit,omitempty" validate:"min=1,max=100"`
}

// DashboardOverview represents the overview statistics
type DashboardOverview struct {
	TotalBusinessTrips     int64   `json:"total_business_trips"`
	DraftBusinessTrips     int64   `json:"draft_business_trips"`
	OngoingBusinessTrips   int64   `json:"ongoing_business_trips"`
	CompletedBusinessTrips int64   `json:"completed_business_trips"`
	CanceledBusinessTrips  int64   `json:"canceled_business_trips"`
	UpcomingBusinessTrips  int64   `json:"upcoming_business_trps"`
	TotalAssignees         int64   `json:"total_assignees"`
	TotalTransactions      int64   `json:"total_transactions"`
	TotalCost              float64 `json:"total_cost"`
	AverageCostPerTrip     float64 `json:"average_cost_per_trip"`
}

// MonthlyStats represents monthly business trip statistics
type MonthlyStats struct {
	Month              string  `json:"month"`
	Year               int     `json:"year"`
	TotalTrips         int64   `json:"total_trips"`
	CompletedTrips     int64   `json:"completed_trips"`
	TotalCost          float64 `json:"total_cost"`
	AverageCostPerTrip float64 `json:"average_cost_per_trip"`
	TopDestination     string  `json:"top_destination"`
}

// DestinationStats represents statistics by destination
type DestinationStats struct {
	Destination        string  `json:"destination"`
	TotalTrips         int64   `json:"total_trips"`
	CompletedTrips     int64   `json:"completed_trips"`
	TotalCost          float64 `json:"total_cost"`
	AverageCostPerTrip float64 `json:"average_cost_per_trip"`
	LastTripDate       string  `json:"last_trip_date"`
}

// TransactionTypeStats represents statistics by transaction type
type TransactionTypeStats struct {
	TransactionType   string  `json:"transaction_type"`
	TotalTransactions int64   `json:"total_transactions"`
	TotalAmount       float64 `json:"total_amount"`
	AverageAmount     float64 `json:"average_amount"`
	Percentage        float64 `json:"percentage"`
}

// RecentBusinessTrip represents a recent business trip for dashboard
type RecentBusinessTrip struct {
	ID                 string  `json:"id"`
	BusinessTripNumber string  `json:"business_trip_number"`
	ActivityPurpose    string  `json:"activity_purpose"`
	DestinationCity    string  `json:"destination_city"`
	StartDate          string  `json:"start_date"`
	EndDate            string  `json:"end_date"`
	Status             string  `json:"status"`
	AssigneeCount      int64   `json:"assignee_count"`
	TotalCost          float64 `json:"total_cost"`
}

// GetDashboardResponse represents the dashboard response
type GetDashboardResponse struct {
	Overview             DashboardOverview      `json:"overview"`
	MonthlyStats         []MonthlyStats         `json:"monthly_stats"`
	DestinationStats     []DestinationStats     `json:"destination_stats"`
	TransactionTypeStats []TransactionTypeStats `json:"transaction_type_stats"`
	RecentBusinessTrips  []RecentBusinessTrip   `json:"recent_business_trips"`
}

// Execute retrieves dashboard data
func (uc *GetDashboardUseCase) Execute(ctx context.Context, req GetDashboardRequest) (*GetDashboardResponse, error) {
	// Get overview statistics using existing repository methods
	overview, err := uc.getOverview(ctx, req)
	if err != nil {
		return nil, err
	}

	// Get monthly statistics
	monthlyStats, err := uc.getMonthlyStats(ctx, req)
	if err != nil {
		return nil, err
	}

	// Get destination statistics
	destinationStats, err := uc.getDestinationStats(ctx, req)
	if err != nil {
		return nil, err
	}

	// Get transaction type statistics
	transactionTypeStats, err := uc.getTransactionTypeStats(ctx, req)
	if err != nil {
		return nil, err
	}

	// Get recent business trips
	recentTrips, err := uc.getRecentBusinessTrips(ctx, req)
	if err != nil {
		return nil, err
	}

	return &GetDashboardResponse{
		Overview:             *overview,
		MonthlyStats:         monthlyStats,
		DestinationStats:     destinationStats,
		TransactionTypeStats: transactionTypeStats,
		RecentBusinessTrips:  recentTrips,
	}, nil
}

// getOverview retrieves overall dashboard statistics
func (uc *GetDashboardUseCase) getOverview(ctx context.Context, req GetDashboardRequest) (*DashboardOverview, error) {
	// Get status counts using repository method
	statusCounts, err := uc.businessTripRepo.GetStatusCounts(ctx, req.StartDate, req.EndDate, req.Destination)
	if err != nil {
		return nil, err
	}

	// Get total cost using repository method
	totalCost, err := uc.businessTripRepo.GetTotalCost(ctx, req.StartDate, req.EndDate, req.Destination)
	if err != nil {
		return nil, err
	}

	// Get total assignees using repository method
	totalAssignees, err := uc.businessTripRepo.GetTotalCount(ctx, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	// Get upcoming trips count
	upcomingCount, err := uc.businessTripRepo.GetUpcomingCount(ctx)
	if err != nil {
		return nil, err
	}

	// Get total transactions count
	var totalTransactions int64 = 0
	if transactionTypeStats, err := uc.transactionRepo.GetTypeStats(ctx, req.StartDate, req.EndDate); err == nil {
		for _, stat := range transactionTypeStats {
			totalTransactions += stat.TotalTransactions
		}
	}

	// Calculate average cost safely to avoid divide by zero
	var averageCost float64
	if statusCounts.Total > 0 {
		averageCost = totalCost / float64(statusCounts.Total)
	}

	return &DashboardOverview{
		TotalBusinessTrips:     statusCounts.Total,
		DraftBusinessTrips:     statusCounts.Draft,
		OngoingBusinessTrips:   statusCounts.Ongoing,
		CompletedBusinessTrips: statusCounts.Completed,
		CanceledBusinessTrips:  statusCounts.Canceled,
		UpcomingBusinessTrips:  upcomingCount,
		TotalAssignees:         totalAssignees,
		TotalTransactions:      totalTransactions,
		TotalCost:              totalCost,
		AverageCostPerTrip:     averageCost,
	}, nil
}

// getMonthlyStats retrieves monthly business trip statistics
func (uc *GetDashboardUseCase) getMonthlyStats(ctx context.Context, req GetDashboardRequest) ([]MonthlyStats, error) {
	// Get last 12 months of data
	now := time.Now()
	twelveMonthsAgo := now.AddDate(-12, 0, 0)

	monthlyData, err := uc.businessTripRepo.GetMonthlyStats(ctx, twelveMonthsAgo, now, req.Destination)
	if err != nil {
		return nil, err
	}

	var monthlyStats []MonthlyStats
	for _, data := range monthlyData {
		var averageCost float64
		if data.TotalTrips > 0 {
			averageCost = data.TotalCost / float64(data.TotalTrips)
		}

		monthlyStats = append(monthlyStats, MonthlyStats{
			Month:              data.Month,
			Year:               data.Year,
			TotalTrips:         data.TotalTrips,
			CompletedTrips:     data.CompletedTrips,
			TotalCost:          data.TotalCost,
			AverageCostPerTrip: averageCost,
			TopDestination:     data.TopDestination,
		})
	}

	return monthlyStats, nil
}

// getDestinationStats retrieves statistics by destination
func (uc *GetDashboardUseCase) getDestinationStats(ctx context.Context, req GetDashboardRequest) ([]DestinationStats, error) {
	destinationData, err := uc.businessTripRepo.GetDestinationStats(ctx, req.StartDate, req.EndDate, req.Destination)
	if err != nil {
		return nil, err
	}

	var destinationStats []DestinationStats
	for _, data := range destinationData {
		var averageCost float64
		if data.TotalTrips > 0 {
			averageCost = data.TotalCost / float64(data.TotalTrips)
		}

		destinationStats = append(destinationStats, DestinationStats{
			Destination:        data.Destination,
			TotalTrips:         data.TotalTrips,
			CompletedTrips:     data.CompletedTrips,
			TotalCost:          data.TotalCost,
			AverageCostPerTrip: averageCost,
			LastTripDate:       data.LastTripDate.Format("2006-01-02"),
		})
	}

	return destinationStats, nil
}

// getTransactionTypeStats retrieves statistics by transaction type
func (uc *GetDashboardUseCase) getTransactionTypeStats(ctx context.Context, req GetDashboardRequest) ([]TransactionTypeStats, error) {
	transactionData, err := uc.transactionRepo.GetTypeStats(ctx, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	// Calculate grand total for percentage calculation
	var grandTotal float64
	for _, data := range transactionData {
		grandTotal += data.TotalAmount
	}

	var transactionTypeStats []TransactionTypeStats
	for _, data := range transactionData {
		var percentage float64
		if grandTotal > 0 {
			percentage = (data.TotalAmount / grandTotal) * 100
		}

		transactionTypeStats = append(transactionTypeStats, TransactionTypeStats{
			TransactionType:   data.TransactionType,
			TotalTransactions: data.TotalTransactions,
			TotalAmount:       data.TotalAmount,
			AverageAmount:     data.AverageAmount,
			Percentage:        percentage,
		})
	}

	return transactionTypeStats, nil
}

// getRecentBusinessTrips retrieves recent business trips
func (uc *GetDashboardUseCase) getRecentBusinessTrips(ctx context.Context, req GetDashboardRequest) ([]RecentBusinessTrip, error) {
	// Set default limit if not provided
	limit := req.Limit
	if limit == 0 {
		limit = 10
	}

	// Get recent business trips with summary data
	recentTrips, err := uc.businessTripRepo.GetRecentWithSummary(ctx, limit)
	if err != nil {
		return nil, err
	}

	var trips []RecentBusinessTrip
	for _, trip := range recentTrips {
		trips = append(trips, RecentBusinessTrip{
			ID:                 trip.ID.String(),
			BusinessTripNumber: trip.BusinessTripNumber,
			ActivityPurpose:    trip.ActivityPurpose,
			DestinationCity:    trip.DestinationCity,
			StartDate:          trip.StartDate.Format("2006-01-02"),
			EndDate:            trip.EndDate.Format("2006-01-02"),
			Status:             string(trip.Status),
			AssigneeCount:      trip.AssigneeCount,
			TotalCost:          trip.TotalCost,
		})
	}

	return trips, nil
}
