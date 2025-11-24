package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// ClientInterface defines the interface for user client to allow for mocking and dependency injection
type ClientInterface interface {
	GetUsersByEmployeeIDs(ctx context.Context, employeeIDs []string) (*UserAPIResponse, error)
	GetSingleUserByEmployeeID(ctx context.Context, employeeID string) (*User, error)
}

// External user API client for fetching user data
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new user API client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		apiKey:  "", // Will be set separately
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetAPIKey sets the API key for the client
func (c *Client) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
}

// NewClientWithAPIKey creates a new user API client with API key
func NewClientWithAPIKey(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// User represents a user from the external API
type User struct {
	ID           string       `json:"id"`
	EmployeeID   string       `json:"employee_id"`
	Username     string       `json:"username"`
	FirstName    string       `json:"first_name"`
	LastName     string       `json:"last_name"`
	Email        *string      `json:"email"`
	PhoneNumber  string       `json:"phone_number"`
	IsActive     bool         `json:"is_active"`
	Organization Organization `json:"organization"`
	Roles        []Role       `json:"roles"`
	CreatedAt    string       `json:"created_at"`
	UpdatedAt    string       `json:"updated_at"`
}

// Organization represents user organization
type Organization struct {
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

// GetUsersByEmployeeIDs fetches users by employee IDs (NIP)
func (c *Client) GetUsersByEmployeeIDs(ctx context.Context, employeeIDs []string) (*UserAPIResponse, error) {
	if len(employeeIDs) == 0 {
		return &UserAPIResponse{Data: []User{}}, nil
	}

	// Build the query string with employee_id parameter (employee_id in external API = NIP)
	employeeIDParam := "in " + strings.Join(employeeIDs, ",")
	url := fmt.Sprintf("%s/users?page=1&limit=%d&employee_id=%s,",
		c.baseURL,
		len(employeeIDs),
		employeeIDParam)

	log.Println(url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call user API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var apiResponse UserAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse user API response: %w", err)
	}

	log.Println(apiResponse)

	return &apiResponse, nil
}

// GetSingleUserByEmployeeID fetches a single user by employee ID
func (c *Client) GetSingleUserByEmployeeID(ctx context.Context, employeeID string) (*User, error) {
	response, err := c.GetUsersByEmployeeIDs(ctx, []string{employeeID})
	if err != nil {
		return nil, err
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("user not found for employee ID: %s", employeeID)
	}

	return &response.Data[0], nil
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

// MockClient is a mock implementation for testing
type MockClient struct {
	users map[string]*User
}

// NewMockClient creates a new mock client with sample data
func NewMockClient() *MockClient {
	return &MockClient{
		users: make(map[string]*User),
	}
}

// AddMockUser adds a mock user to the client
func (m *MockClient) AddMockUser(employeeID string, user *User) {
	m.users[employeeID] = user
}

// GetUsersByEmployeeIDs mock implementation
func (m *MockClient) GetUsersByEmployeeIDs(ctx context.Context, employeeIDs []string) (*UserAPIResponse, error) {
	var users []User
	for _, id := range employeeIDs {
		if user, exists := m.users[id]; exists {
			users = append(users, *user)
		}
	}

	return &UserAPIResponse{
		Data:       users,
		Page:       1,
		Limit:      len(employeeIDs),
		TotalItems: len(users),
		TotalPages: 1,
	}, nil
}

// GetSingleUserByEmployeeID mock implementation
func (m *MockClient) GetSingleUserByEmployeeID(ctx context.Context, employeeID string) (*User, error) {
	user, exists := m.users[employeeID]
	if !exists {
		return nil, fmt.Errorf("user not found for employee ID: %s", employeeID)
	}
	return user, nil
}
