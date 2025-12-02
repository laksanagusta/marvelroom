package work_paper

// OrganizationResponse represents organization data in the response
type OrganizationResponse struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Address *string `json:"address,omitempty"`
	Type    string  `json:"type"`
}
