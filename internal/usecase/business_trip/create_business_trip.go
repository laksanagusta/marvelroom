package business_trip

import (
	"context"
	"encoding/json"
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
	historyUseCase   *RecordHistoryUseCase
}

func NewCreateBusinessTripUseCase(
	businessTripRepo repository.BusinessTripRepository,
	assigneeRepo repository.AssigneeRepository,
	transactionRepo repository.BusinessTripTransactionRepository,
	userService *service.UserService,
	db database.DB,
	historyUseCase *RecordHistoryUseCase,
) *CreateBusinessTripUseCase {
	return &CreateBusinessTripUseCase{
		businessTripRepo: businessTripRepo,
		assigneeRepo:     assigneeRepo,
		transactionRepo:  transactionRepo,
		userService:      userService,
		db:               db,
		historyUseCase:   historyUseCase,
	}
}

func (uc *CreateBusinessTripUseCase) Execute(ctx context.Context, req BusinessTripRequest, authenticatedUser entity.AuthenticatedUser) (*BusinessTripResponse, error) {
	var employeeNumbers []string
	for _, assigneeReq := range req.Assignees {
		if assigneeReq.EmployeeNumber != "" {
			employeeNumbers = append(employeeNumbers, assigneeReq.EmployeeNumber)
		}
	}

	jsn, _ := json.Marshal(req)
	fmt.Println("req", string(jsn))

	userDataMap, err := uc.userService.GetUserDataByEmployeeIDs(ctx, employeeNumbers)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user data: %w", err)
	}

	var invalidEmployees []InvalidEmployeeError
	for i, assigneeReq := range req.Assignees {
		if assigneeReq.EmployeeNumber != "" {
			if _, exists := userDataMap[assigneeReq.EmployeeNumber]; !exists {
				invalidEmployees = append(invalidEmployees, InvalidEmployeeError{
					Index:          i,
					Field:          "employee_number",
					EmployeeNumber: assigneeReq.EmployeeNumber,
					Name:           assigneeReq.Name,
					Message:        fmt.Sprintf("Employee number '%s' not found in identity service", assigneeReq.EmployeeNumber),
				})
			}
		}
	}

	if len(invalidEmployees) > 0 {
		return nil, &EmployeeValidationError{
			Message:          "Some employee numbers were not found in identity service",
			InvalidEmployees: invalidEmployees,
		}
	}

	bt, err := req.ToEntity(authenticatedUser.Organization.ID)
	if err != nil {
		return nil, err
	}

	var completeBusinessTrip *entity.BusinessTrip

	err = uc.db.WithTransaction(ctx, func(ctx context.Context, tx database.DBTx) error {
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

		if uc.historyUseCase != nil {
			err = uc.historyUseCase.ExecuteWithTx(ctx, RecordHistoryInput{
				BusinessTripID: businessTrip.ID,
				ChangeType:     entity.HistoryChangeTypeStatusChange,
				FieldName:      "status",
				OldValue:       "",
				NewValue:       string(businessTrip.Status),
				ChangedBy:      authenticatedUser.GetFullName(),
			}, tx)
			if err != nil {
				return fmt.Errorf("failed to record history: %w", err)
			}
		}

		for _, assignee := range businessTrip.Assignees {
			assignee.BusinessTripID = businessTrip.ID

			employeeNumber := assignee.EmployeeNumber
			if employeeNumber == "" {
				employeeNumber = assignee.EmployeeID // fallback
			}

			if userData, exists := userDataMap[employeeNumber]; exists {
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
				jsn, _ := json.Marshal(transaction)
				fmt.Println("transaction", string(jsn))
				_, err := transactionRepoWithTx.CreateTransaction(ctx, transaction)
				if err != nil {
					return err
				}
			}
		}

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

	completeBusinessTrip, err = uc.businessTripRepo.GetByID(ctx, bt.ID)
	if err != nil {
		return nil, err
	}

	return FromEntity(completeBusinessTrip), nil
}
