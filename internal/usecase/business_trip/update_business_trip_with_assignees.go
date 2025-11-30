package business_trip

import (
	"context"
	"fmt"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/internal/domain/service"
	"sandbox/pkg/database"
)

type UpdateBusinessTripWithAssigneesUseCase struct {
	businessTripRepo repository.BusinessTripRepository
	assigneeRepo     repository.AssigneeRepository
	transactionRepo  repository.BusinessTripTransactionRepository
	userService      *service.UserService
	db               database.DB
}

func NewUpdateBusinessTripWithAssigneesUseCase(businessTripRepo repository.BusinessTripRepository, assigneeRepo repository.AssigneeRepository, transactionRepo repository.BusinessTripTransactionRepository, userService *service.UserService, db database.DB) *UpdateBusinessTripWithAssigneesUseCase {
	return &UpdateBusinessTripWithAssigneesUseCase{
		businessTripRepo: businessTripRepo,
		assigneeRepo:     assigneeRepo,
		transactionRepo:  transactionRepo,
		userService:      userService,
		db:               db,
	}
}

func (uc *UpdateBusinessTripWithAssigneesUseCase) Execute(ctx context.Context, req UpdateBusinessTripWithAssigneesRequest) (*BusinessTripResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Extract employee numbers from request to fetch user data
	var employeeNumbers []string
	for _, assigneeReq := range req.Assignees {
		if assigneeReq.EmployeeNumber != "" {
			employeeNumbers = append(employeeNumbers, assigneeReq.EmployeeNumber)
		} else if assigneeReq.EmployeeID != "" {
			// Fallback to employee_id if employee_number is not provided
			employeeNumbers = append(employeeNumbers, assigneeReq.EmployeeID)
		}
	}

	// Fetch user data from external API
	userDataMap, err := uc.userService.GetUserDataByEmployeeIDs(ctx, employeeNumbers)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user data: %w", err)
	}

	bt, err := req.ToEntity(req.BusinessTripID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert request to entity: %w", err)
	}

	var result *entity.BusinessTrip
	err = uc.db.WithTransaction(ctx, func(ctx context.Context, tx database.DBTx) error {
		repoWithTx := uc.businessTripRepo.(interface {
			WithTransaction(database.DBTx) repository.BusinessTripRepository
		}).WithTransaction(tx)

		assigneeRepoWithTx := uc.assigneeRepo.(interface {
			WithTransaction(database.DBTx) repository.AssigneeRepository
		}).WithTransaction(tx)

		transactionRepoWithTx := uc.transactionRepo.(interface {
			WithTransaction(database.DBTx) repository.BusinessTripTransactionRepository
		}).WithTransaction(tx)

		_, err = repoWithTx.Update(ctx, bt)
		if err != nil {
			return fmt.Errorf("failed to update business trip: %w", err)
		}

		existingAssignees, err := assigneeRepoWithTx.GetAssigneesByBusinessTripIDWithoutTransactions(ctx, req.BusinessTripID)
		if err != nil {
			return fmt.Errorf("failed to get existing assignees: %w", err)
		}

		if len(existingAssignees) > 0 {
			assigneeIDs := make([]string, len(existingAssignees))
			for i, assignee := range existingAssignees {
				assigneeIDs[i] = assignee.ID
			}
			err = repoWithTx.DeleteTransactionsByAssigneeIDs(ctx, assigneeIDs)
			if err != nil {
				return fmt.Errorf("failed to delete existing transactions: %w", err)
			}
		}

		err = assigneeRepoWithTx.DeleteAssigneesByBusinessTripID(ctx, req.BusinessTripID)
		if err != nil {
			return fmt.Errorf("failed to delete existing assignees: %w", err)
		}

		for _, assignee := range bt.Assignees {
			assignee.BusinessTripID = req.BusinessTripID

			// Find employee number to fetch user data
			employeeNumber := assignee.EmployeeNumber
			if employeeNumber == "" {
				employeeNumber = assignee.EmployeeID // fallback
			}

			// Fetch user data from external API
			if userData, exists := userDataMap[employeeNumber]; exists {
				// Update assignee with data from external API
				assignee.EmployeeID = userData.EmployeeID         // external API user ID
				assignee.EmployeeName = userData.Name             // full name from API
				assignee.EmployeeNumber = userData.EmployeeNumber // NIP from API
			}

			createdAssignee, err := assigneeRepoWithTx.Create(ctx, assignee)
			if err != nil {
				return fmt.Errorf("failed to create assignee %s: %w", assignee.Name, err)
			}

			for _, transaction := range assignee.Transactions {
				transaction.AssigneeID = createdAssignee.ID

				_, err = transactionRepoWithTx.CreateTransaction(ctx, transaction)
				if err != nil {
					return fmt.Errorf("failed to create transaction %s for assignee %s: %w", transaction.Name, assignee.Name, err)
				}
			}
		}

		// Construct the complete business trip entity from the data we already have
		// This avoids the problematic GetByID query that was causing PostgreSQL protocol issues
		assignees, err := assigneeRepoWithTx.GetAssigneesByBusinessTripIDWithoutTransactions(ctx, req.BusinessTripID)
		if err != nil {
			return fmt.Errorf("failed to get updated assignees: %w", err)
		}

		// Load transactions for each assignee using the transaction-aware repository
		for _, assignee := range assignees {
			transactions, err := transactionRepoWithTx.GetTransactionsByAssigneeID(ctx, assignee.ID)
			if err != nil {
				return fmt.Errorf("failed to get transactions for assignee %s: %w", assignee.ID, err)
			}
			assignee.Transactions = transactions
		}

		bt.Assignees = assignees
		result = bt
		return nil
	})
	if err != nil {
		return nil, err
	}

	return FromEntity(result), nil
}
