package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"sandbox/internal/domain/entity"
)

// User represents a user from external API
type User struct {
	ID           string           `json:"id"`
	EmployeeID   string           `json:"employee_id"`
	Username     string           `json:"username"`
	FirstName    string           `json:"first_name"`
	LastName     string           `json:"last_name"`
	Email        *string          `json:"email"`
	PhoneNumber  string           `json:"phone_number"`
	IsActive     bool             `json:"is_active"`
	Organization UserOrganization `json:"organization"`
	Roles        []Role           `json:"roles"`
	CreatedAt    string           `json:"created_at"`
	UpdatedAt    string           `json:"updated_at"`
}

// UserOrganization represents user organization (different from domain entity.Organization)
type UserOrganization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// Role represents user role
type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UserAPIResponse represents the API response structure
type UserAPIResponse struct {
	Data       []User `json:"data"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalItems int    `json:"total_items"`
	TotalPages int    `json:"total_pages"`
}

// CreateUserData represents the data needed to create a user in our system
type CreateUserData struct {
	EmployeeID     string
	EmployeeNumber string
	Name           string
	FirstName      string
	LastName       string
	Email          string
	PhoneNumber    string
	Organization   string
}

// ExtractCreateUserData extracts relevant user data for our system
func (u *User) ExtractCreateUserData() *CreateUserData {
	name := u.FirstName
	if u.LastName != "" {
		if name != "" {
			name += " " + u.LastName
		} else {
			name = u.LastName
		}
	}

	email := ""
	if u.Email != nil {
		email = *u.Email
	}

	return &CreateUserData{
		EmployeeID:     u.ID,         // external API response.id becomes our employee_id
		EmployeeNumber: u.EmployeeID, // external API response.employee_id (NIP) becomes our employee_number
		Name:           name,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Email:          email,
		PhoneNumber:    u.PhoneNumber,
		Organization:   u.Organization.Name,
	}
}

// IdentityService handles communication with unified identity API (users + organizations)
type IdentityService struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// IdentityServiceInterface defines the interface for identity service
type IdentityServiceInterface interface {
	GetUsersByEmployeeIDs(ctx context.Context, employeeIDs []string) (*UserAPIResponse, error)
	GetSingleUserByEmployeeID(ctx context.Context, employeeID string) (*User, error)
	GetOrganizations(ctx context.Context, page, limit int, sort string) (*entity.OrganizationListResponse, error)
	GetOrganizationByID(ctx context.Context, id string) (*entity.Organization, error)
}

// NewIdentityService creates a new identity service
func NewIdentityService(baseURL string) *IdentityService {
	return &IdentityService{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewIdentityServiceWithAPIKey creates a new identity service with API key
func NewIdentityServiceWithAPIKey(baseURL, apiKey string) *IdentityService {
	return &IdentityService{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetAPIKey sets the API key for the service
func (s *IdentityService) SetAPIKey(apiKey string) {
	s.apiKey = apiKey
}

// GetUsersByEmployeeIDs fetches users by employee IDs
func (s *IdentityService) GetUsersByEmployeeIDs(ctx context.Context, employeeIDs []string) (*UserAPIResponse, error) {
	if len(employeeIDs) == 0 {
		return &UserAPIResponse{Data: []User{}}, nil
	}

	employeeIDParam := "in " + strings.Join(employeeIDs, ",")
	url := fmt.Sprintf("%s/users?page=1&limit=%d&employee_id=%s,",
		s.baseURL,
		len(employeeIDs),
		employeeIDParam)

	log.Println("url", url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if s.apiKey != "" {
		req.Header.Set("X-API-Key", s.apiKey)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call identity API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("identity API error (status %d)", resp.StatusCode)
	}

	var apiResponse UserAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse user API response: %w", err)
	}

	return &apiResponse, nil
}

// GetSingleUserByEmployeeID fetches a single user by employee ID
func (s *IdentityService) GetSingleUserByEmployeeID(ctx context.Context, employeeID string) (*User, error) {
	response, err := s.GetUsersByEmployeeIDs(ctx, []string{employeeID})
	if err != nil {
		return nil, err
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("user not found for employee ID: %s", employeeID)
	}

	return &response.Data[0], nil
}

// GetOrganizations fetches organizations from identity API
func (s *IdentityService) GetOrganizations(ctx context.Context, page, limit int, sort string) (*entity.OrganizationListResponse, error) {
	url := fmt.Sprintf("%s/api/v1/organizations?page=%d&limit=%d&sort=%s", s.baseURL, page, limit, sort)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if s.apiKey != "" {
		req.Header.Set("X-API-Key", s.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call identity API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("organizations API returned status: %d", resp.StatusCode)
	}

	var response entity.OrganizationListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// GetOrganizationByID fetches a specific organization by ID
func (s *IdentityService) GetOrganizationByID(ctx context.Context, id string) (*entity.Organization, error) {
	url := fmt.Sprintf("%s/organizations/%s", s.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if s.apiKey != "" {
		req.Header.Set("X-API-Key", s.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call identity API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, entity.ErrOrganizationNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("organizations API returned status: %d", resp.StatusCode)
	}

	// Parse the response with expected structure
	var apiResponse struct {
		Data entity.Organization `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &apiResponse.Data, nil
}
