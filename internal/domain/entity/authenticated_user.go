package entity

import (
	"github.com/google/uuid"
)

// Role represents a user role
type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UserOrganization represents simplified user organization info
type UserOrganization struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Type string    `json:"type"`
}

type AuthenticatedUser struct {
	ID           string           `json:"id"`
	EmployeeID   string           `json:"employee_id"`
	Username     string           `json:"username"`
	FirstName    string           `json:"first_name"`
	LastName     string           `json:"last_name"`
	PhoneNumber  string           `json:"phone_number"`
	Roles        []Role           `json:"roles"`
	Organization UserOrganization `json:"organization"`
}

// GetFullName returns the user's full name
func (u *AuthenticatedUser) GetFullName() string {
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	if u.FirstName != "" {
		return u.FirstName
	}
	if u.LastName != "" {
		return u.LastName
	}
	return u.Username
}

// GetPrimaryRole returns the user's primary role name (first role if multiple)
func (u *AuthenticatedUser) GetPrimaryRole() string {
	if len(u.Roles) > 0 {
		return u.Roles[0].Name
	}
	return ""
}

// HasRole checks if user has a specific role by name
func (u *AuthenticatedUser) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}
