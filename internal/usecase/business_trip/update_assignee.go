package business_trip

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/invopop/validation"

	"sandbox/internal/domain/repository"
	"sandbox/internal/domain/service"
	"sandbox/internal/infrastructure/user"
)

type UpdateAssigneeUseCase struct {
	businessTripRepo repository.BusinessTripRepository
	userService     *service.UserService
}

func NewUpdateAssigneeUseCase(businessTripRepo repository.BusinessTripRepository, userClient user.ClientInterface) *UpdateAssigneeUseCase {
	userService := service.NewUserService(userClient)
	return &UpdateAssigneeUseCase{
		businessTripRepo: businessTripRepo,
		userService:     userService,
	}
}

type UpdateAssigneeRequest struct {
	BusinessTripID string `params:"businessTripId" json:"businessTripId"`
	AssigneeID     string `params:"assigneeId" json:"assigneeId"`
	Name           string `json:"name"`
	SPDNumber      string `json:"spdNumber"`
	EmployeeID     string `json:"employeeId"`
	EmployeeName   string `json:"employeeName"`
	EmployeeNumber string `json:"employeeNumber"`
	Position       string `json:"position"`
	Rank           string `json:"rank"`
}

func (r UpdateAssigneeRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.BusinessTripID, validation.Required),
		validation.Field(&r.AssigneeID, validation.Required),
		validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.SPDNumber, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.EmployeeNumber, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Position, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Rank, validation.Required, validation.Length(1, 100)),
	)
}

type UpdateAssigneeResponse struct {
	ID             string `json:"id"`
	BusinessTripID string `json:"businessTripId"`
	Name           string `json:"name"`
	SPDNumber      string `json:"spdNumber"`
	EmployeeID     string `json:"employeeId"`
	EmployeeName   string `json:"employeeName"`
	EmployeeNumber string `json:"employeeNumber"`
	Position       string `json:"position"`
	Rank           string `json:"rank"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

func (uc *UpdateAssigneeUseCase) Execute(ctx context.Context, req UpdateAssigneeRequest) (*UpdateAssigneeResponse, error) {
	businessTrip, err := uc.businessTripRepo.GetByID(ctx, req.BusinessTripID)
	if err != nil {
		return nil, fmt.Errorf("failed to get business trip: %w", err)
	}
	if businessTrip == nil {
		return nil, fmt.Errorf("business trip not found")
	}

	assignee, err := uc.businessTripRepo.GetAssigneeByID(ctx, req.AssigneeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignee: %w", err)
	}
	if assignee == nil {
		return nil, fmt.Errorf("assignee not found")
	}

	if assignee.BusinessTripID != req.BusinessTripID {
		return nil, fmt.Errorf("assignee does not belong to the specified business trip")
	}

		assignees, err := uc.businessTripRepo.GetAssigneesByBusinessTripID(ctx, assignee.BusinessTripID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignees: %w", err)
	}

	for _, existingAssignee := range assignees {
		if existingAssignee.ID != req.AssigneeID && strings.EqualFold(existingAssignee.SPDNumber, req.SPDNumber) {
			return nil, fmt.Errorf("SPD number %s already exists for this business trip", req.SPDNumber)
		}
	}

	// Fetch user data from external API if employee_number is provided or if employee_id has changed
	var employeeNumber string
	if req.EmployeeNumber != "" {
		employeeNumber = req.EmployeeNumber
	} else if req.EmployeeID != assignee.EmployeeID {
		// If employee_id changed, use new employee_id to fetch data
		employeeNumber = req.EmployeeID
	} else {
		// Use existing employee number
		employeeNumber = assignee.EmployeeNumber
	}

	userData, err := uc.userService.GetSingleUserDataByEmployeeID(ctx, employeeNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user data for employee number %s: %w", employeeNumber, err)
	}

	// Update assignee details with data from external API
	assignee.Name = userData.Name
	assignee.SPDNumber = strings.TrimSpace(req.SPDNumber)
	assignee.EmployeeID = userData.EmployeeID
	assignee.EmployeeName = userData.Name
	assignee.EmployeeNumber = userData.EmployeeNumber
	assignee.Position = strings.TrimSpace(req.Position) // Keep position from request as it might be specific to the trip
	assignee.Rank = strings.TrimSpace(req.Rank)         // Keep rank from request as it might be specific to the trip

	// Save updated assignee
	updatedAssignee, err := uc.businessTripRepo.UpdateAssignee(ctx, assignee)
	if err != nil {
		return nil, fmt.Errorf("failed to update assignee: %w", err)
	}

	return &UpdateAssigneeResponse{
		ID:             updatedAssignee.ID,
		BusinessTripID: updatedAssignee.BusinessTripID,
		Name:           updatedAssignee.Name,
		SPDNumber:      updatedAssignee.SPDNumber,
		EmployeeID:     updatedAssignee.EmployeeID,
		EmployeeName:   updatedAssignee.EmployeeName,
		EmployeeNumber: updatedAssignee.EmployeeNumber,
		Position:       updatedAssignee.Position,
		Rank:           updatedAssignee.Rank,
		CreatedAt:      updatedAssignee.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      updatedAssignee.UpdatedAt.Format(time.RFC3339),
	}, nil
}
