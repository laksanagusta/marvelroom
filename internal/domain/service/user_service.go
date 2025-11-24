package service

import (
	"context"
	"fmt"

	"sandbox/internal/infrastructure/user"
)

// UserService handles user-related operations including fetching user data from external API
type UserService struct {
	userClient user.ClientInterface
}

// NewUserService creates a new user service
func NewUserService(userClient user.ClientInterface) *UserService {
	return &UserService{
		userClient: userClient,
	}
}

// GetUserDataByEmployeeIDs fetches user data for multiple employee IDs (NIP) from external API
func (s *UserService) GetUserDataByEmployeeIDs(ctx context.Context, employeeNumbers []string) (map[string]*user.CreateUserData, error) {
	if len(employeeNumbers) == 0 {
		return make(map[string]*user.CreateUserData), nil
	}

	// Fetch users from external API using NIP numbers
	response, err := s.userClient.GetUsersByEmployeeIDs(ctx, employeeNumbers)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users from external API: %w", err)
	}

	// Convert to map for easy lookup using employee_number (NIP) as key
	userDataMap := make(map[string]*user.CreateUserData)
	for _, apiUser := range response.Data {
		userData := apiUser.ExtractCreateUserData()
		userDataMap[userData.EmployeeNumber] = userData // Key by NIP (employee_number)
	}

	return userDataMap, nil
}

// GetSingleUserDataByEmployeeID fetches user data for a single employee ID (NIP) from external API
func (s *UserService) GetSingleUserDataByEmployeeID(ctx context.Context, employeeNumber string) (*user.CreateUserData, error) {
	apiUser, err := s.userClient.GetSingleUserByEmployeeID(ctx, employeeNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user from external API: %w", err)
	}

	userData := apiUser.ExtractCreateUserData()
	return userData, nil
}

// ValidateEmployeeNumbers checks if all employee numbers (NIP) exist in the external user service
func (s *UserService) ValidateEmployeeNumbers(ctx context.Context, employeeNumbers []string) (map[string]bool, error) {
	if len(employeeNumbers) == 0 {
		return make(map[string]bool), nil
	}

	response, err := s.userClient.GetUsersByEmployeeIDs(ctx, employeeNumbers)
	if err != nil {
		return nil, fmt.Errorf("failed to validate employee numbers: %w", err)
	}

	validationMap := make(map[string]bool)

	// Initialize all as false
	for _, number := range employeeNumbers {
		validationMap[number] = false
	}

	// Mark found numbers as true (using employee_id from external API which is NIP)
	for _, apiUser := range response.Data {
		validationMap[apiUser.EmployeeID] = true
	}

	return validationMap, nil
}

// ValidateEmployeeIDs is a legacy method - use ValidateEmployeeNumbers instead
func (s *UserService) ValidateEmployeeIDs(ctx context.Context, employeeIDs []string) (map[string]bool, error) {
	return s.ValidateEmployeeNumbers(ctx, employeeIDs)
}