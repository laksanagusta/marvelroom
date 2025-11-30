package entity

import (
	"time"

	"github.com/google/uuid"
)

// Organization represents an organizational unit from the external API
type Organization struct {
	ID            uuid.UUID      `json:"id"`
	Name          string         `json:"name"`
	Address       *string        `json:"address"`
	Type          string         `json:"type"`
	Organizations []Organization `json:"organizations,omitempty"` // Nested organizations
	CreatedAt     time.Time      `json:"created_at"`
	CreatedBy     string         `json:"created_by"`
}

// OrganizationListResponse represents the response from organizations API
type OrganizationListResponse struct {
	Data     []Organization `json:"data"`
	Metadata Metadata       `json:"metadata"`
}

// Metadata represents pagination metadata from organizations API
type Metadata struct {
	Count       int `json:"count"`
	TotalCount  int `json:"total_count"`
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
}
