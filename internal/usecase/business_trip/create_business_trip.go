package business_trip

import (
	"context"
	"fmt"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/internal/domain/service"
	"sandbox/pkg/database"
)

type CreateBusinessTripUseCase struct {
	businessTripRepo repository.BusinessTripRepository
	assigneeRepo     repository.AssigneeRepository
	transactionRepo  repository.BusinessTripTransactionRepository
	userService      *service.UserService
	db               database.DB
}

func NewCreateBusinessTripUseCase(businessTripRepo repository.BusinessTripRepository, assigneeRepo repository.AssigneeRepository, transactionRepo repository.BusinessTripTransactionRepository, userService *service.UserService, db database.DB) *CreateBusinessTripUseCase {
	return &CreateBusinessTripUseCase{
		businessTripRepo: businessTripRepo,
		assigneeRepo:     assigneeRepo,
		transactionRepo:  transactionRepo,
		userService:      userService,
		db:               db,
	}
}

func (uc *CreateBusinessTripUseCase) Execute(ctx context.Context, req BusinessTripRequest) (*BusinessTripResponse, error) {
	// Extract employee numbers from request to fetch user data
	var employeeNumbers []string
	for _, assigneeReq := range req.Assignees {
		if assigneeReq.EmployeeNumber != "" {
			employeeNumbers = append(employeeNumbers, assigneeReq.EmployeeNumber)
		}
	}

	// Fetch user data from external API
	userDataMap, err := uc.userService.GetUserDataByEmployeeIDs(ctx, employeeNumbers)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user data: %w", err)
	}

	bt, err := req.ToEntity()
	if err != nil {
		return nil, err
	}

	var completeBusinessTrip *entity.BusinessTrip

	err = uc.db.WithTransaction(ctx, func(ctx context.Context, tx database.DBTx) error {
		// Create transaction-aware repositories
		businessTripRepoWithTx := uc.businessTripRepo.(interface {
			WithTransaction(database.DBTx) repository.BusinessTripRepository
		}).WithTransaction(tx)

		assigneeRepoWithTx := uc.assigneeRepo.(interface {
			WithTransaction(database.DBTx) repository.AssigneeRepository
		}).WithTransaction(tx)

		transactionRepoWithTx := uc.transactionRepo.(interface {
			WithTransaction(database.DBTx) repository.BusinessTripTransactionRepository
		}).WithTransaction(tx)

		businessTrip, err := businessTripRepoWithTx.Create(ctx, bt)
		if err != nil {
			return err
		}

		// Process assignees with external API data
		for _, assignee := range businessTrip.Assignees {
			assignee.BusinessTripID = businessTrip.ID

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
				return err
			}

			for _, transaction := range createdAssignee.Transactions {
				transaction.AssigneeID = createdAssignee.ID
				_, err := transactionRepoWithTx.CreateTransaction(ctx, transaction)
				if err != nil {
					return err
				}
			}
		}

		// Create verificators in database
		for _, verificator := range businessTrip.Verificators {
			verificator.BusinessTripID = businessTrip.ID
			_, err := businessTripRepoWithTx.CreateVerificator(ctx, verificator)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Query complete business trip using main repository (not within transaction)
	// This avoids connection state issues within transaction
	completeBusinessTrip, err = uc.businessTripRepo.GetByID(ctx, bt.ID)
	if err != nil {
		return nil, err
	}

	return FromEntity(completeBusinessTrip), nil
}
