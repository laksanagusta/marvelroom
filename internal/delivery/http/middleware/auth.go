package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"sandbox/internal/domain/entity"
)

// Role represents a user role from identity service
type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Organization represents user organization from identity service
type Organization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// IdentityData represents the data field from identity service response
type IdentityData struct {
	ID            string       `json:"id"`
	EmployeeID    string       `json:"employee_id"`
	Username      string       `json:"username"`
	FirstName     string       `json:"first_name"`
	LastName      string       `json:"last_name"`
	PhoneNumber   string       `json:"phone_number"`
	Roles         []Role       `json:"roles"`
	Permissions   interface{}  `json:"permissions"`
	Organization  Organization `json:"organization"`
	Scopes        interface{}  `json:"scopes"`
}

// IdentityResponse represents the response from identity service
type IdentityResponse struct {
	Data IdentityData `json:"data"`
}

// AuthMiddleware creates a middleware that checks authentication with identity service
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		// Check Bearer token format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format. Expected: Bearer <token>",
			})
		}

		token := tokenParts[1]

		// Call identity service /whoami API
		user, err := callIdentityService(token)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": fmt.Sprintf("Authentication failed: %v", err),
			})
		}

		// Store authenticated user in context locals
		c.Locals("authenticatedUser", user)

		// Continue to next handler
		return c.Next()
	}
}

// callIdentityService calls the identity service to validate token and get user info
func callIdentityService(token string) (*entity.AuthenticatedUser, error) {
	// Create HTTP client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", "http://localhost:5001/api/v1/users/whoami", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set Authorization header
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call identity service: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("identity service returned status: %d", resp.StatusCode)
	}

	// Read response body
	var body bytes.Buffer
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response
	var identityResp IdentityResponse
	if err := json.Unmarshal(body.Bytes(), &identityResp); err != nil {
		return nil, fmt.Errorf("failed to parse identity response: %w", err)
	}

	// Validate required fields
	if identityResp.Data.ID == "" {
		return nil, errors.New("invalid identity response: missing user ID")
	}

	// Convert Identity roles to entity roles
	roles := make([]entity.Role, len(identityResp.Data.Roles))
	for i, role := range identityResp.Data.Roles {
		roles[i] = entity.Role{
			ID:   role.ID,
			Name: role.Name,
		}
	}

	// Parse organization ID as UUID
	orgUUID, err := uuid.Parse(identityResp.Data.Organization.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID format: %w", err)
	}

	// Convert Identity organization to UserOrganization entity
	organization := entity.UserOrganization{
		ID:   orgUUID,
		Name: identityResp.Data.Organization.Name,
		Type: identityResp.Data.Organization.Type,
	}

	// Convert to AuthenticatedUser entity
	user := &entity.AuthenticatedUser{
		ID:           identityResp.Data.ID,
		EmployeeID:   identityResp.Data.EmployeeID,
		Username:     identityResp.Data.Username,
		FirstName:    identityResp.Data.FirstName,
		LastName:     identityResp.Data.LastName,
		PhoneNumber:  identityResp.Data.PhoneNumber,
		Roles:        roles,
		Organization: organization,
	}

	return user, nil
}

// GetAuthenticatedUser retrieves the authenticated user from context
func GetAuthenticatedUser(c *fiber.Ctx) (*entity.AuthenticatedUser, error) {
	user, ok := c.Locals("authenticatedUser").(*entity.AuthenticatedUser)
	if !ok {
		return nil, errors.New("authenticated user not found in context")
	}
	return user, nil
}