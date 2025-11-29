package business_trip

import (
	"context"
	"fmt"
	"time"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/internal/domain/service"
	"sandbox/pkg/database"
)

type AddAssigneeUseCase struct {
	businessTripRepo repository.BusinessTripRepository
	userService     *service.UserService
	db               database.DB
}

func NewAddAssigneeUseCase(businessTripRepo repository.BusinessTripRepository, userService *service.UserService, db database.DB) *AddAssigneeUseCase {
	return &AddAssigneeUseCase{
		businessTripRepo: businessTripRepo,
		userService:     userService,
		db:               db,
	}
}

func (uc *AddAssigneeUseCase) Execute(ctx context.Context, businessTripID string, req *AssigneeRequest) (*AssigneeResponse, error) {
	businessTrip, err := uc.businessTripRepo.GetByID(ctx, businessTripID)
	if err != nil {
		return nil, err
	}
	if businessTrip == nil {
		return nil, entity.ErrBusinessTripNotFound
	}

	// Fetch user data from external API using employee_number (which comes from req.EmployeeNumber)
	// Note: According to the requirement, employee_number is the result from external API where employee_id = employee_number
	var employeeNumber string
	if req.EmployeeNumber != "" {
		employeeNumber = req.EmployeeNumber
	} else {
		// If employee_number is not provided in request, use employee_id as fallback
		employeeNumber = req.EmployeeID
	}

	userData, err := uc.userService.GetSingleUserDataByEmployeeID(ctx, employeeNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user data for employee number %s: %w", employeeNumber, err)
	}

	// Create assignee with data from external API
	assignee := &entity.Assignee{
		Name:           userData.Name,
		SPDNumber:      req.SPDNumber,
		EmployeeID:     userData.EmployeeID,
		EmployeeName:   userData.Name,
		EmployeeNumber: userData.EmployeeNumber,
		Position:       req.Position, // Keep position from request as it might be specific to the trip
		Rank:           req.Rank,     // Keep rank from request as it might be specific to the trip
	}

	for _, txReq := range req.Transactions {
		transaction := &entity.Transaction{
			Name:            txReq.Name,
			Type:            entity.TransactionType(txReq.Type),
			Subtype:         entity.TransactionSubtype(txReq.Subtype),
			Amount:          txReq.Amount,
			TotalNight:      txReq.TotalNight,
			Description:     txReq.Description,
			TransportDetail: txReq.TransportDetail,
		}
		assignee.Transactions = append(assignee.Transactions, transaction)
	}

	// Set business trip ID before creating
	assignee.BusinessTripID = businessTripID

	var createdAssignee *entity.Assignee

	err = uc.db.WithTransaction(ctx, func(ctx context.Context, tx database.DBTx) error {
		// Create transaction-aware repository
		repoWithTx := uc.businessTripRepo.(interface {
			WithTransaction(database.DBTx) repository.BusinessTripRepository
		}).WithTransaction(tx)

		var err error
		createdAssignee, err = repoWithTx.CreateAssignee(ctx, assignee)
		if err != nil {
			return err
		}

		// Create transactions for the assignee
		for _, transaction := range assignee.Transactions {
			transaction.AssigneeID = createdAssignee.ID
			_, err = repoWithTx.CreateTransaction(ctx, transaction)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &AssigneeResponse{
		ID:             createdAssignee.GetID(),
		Name:           createdAssignee.GetName(),
		SPDNumber:      createdAssignee.GetSPDNumber(),
		EmployeeID:     createdAssignee.GetEmployeeID(),
		EmployeeName:   createdAssignee.GetEmployeeName(),
		EmployeeNumber: createdAssignee.GetEmployeeNumber(),
		Position:       createdAssignee.GetPosition(),
		Rank:           createdAssignee.GetRank(),
		TotalCost:      createdAssignee.GetTotalCost(),
		Transactions:   []TransactionResponse{}, // Empty for now
		CreatedAt:      createdAssignee.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      createdAssignee.UpdatedAt.Format(time.RFC3339),
	}, nil
}
