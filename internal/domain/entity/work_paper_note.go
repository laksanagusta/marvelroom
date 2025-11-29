package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// WorkPaperNote represents a note within a work paper
type WorkPaperNote struct {
	ID             uuid.UUID     `db:"id"`
	WorkPaperID    uuid.UUID     `db:"work_paper_id"`
	MasterItemID   uuid.UUID     `db:"master_item_id"`
	GDriveLink     *string       `db:"gdrive_link"`           // Nullable Google Drive link
	IsValid        *bool         `db:"is_valid"`              // Nullable Y/T result from LLM
	Notes          *string       `db:"notes"`                 // Nullable notes from LLM
	LastLLMResponse *LLMResponse `db:"last_llm_response"`    // Nullable raw LLM response
	CreatedAt      time.Time     `db:"created_at"`
	UpdatedAt      time.Time     `db:"updated_at"`
	DeletedAt      *time.Time    `db:"deleted_at"`

	// Relations
	MasterItem   *WorkPaperItem  `db:"-"`
}

// LLMResponse represents the response structure from LLM
type LLMResponse struct {
	Note     string `json:"note"`
	IsValid  bool   `json:"isValid"`
	Model    string `json:"model,omitempty"`
	Usage    *Usage `json:"usage,omitempty"`
}

// Usage represents token usage from LLM API
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Value implements the driver.Valuer interface for LLMResponse
func (r LLMResponse) Value() (driver.Value, error) {
	return json.Marshal(r)
}

// Scan implements the sql.Scanner interface for LLMResponse
func (r *LLMResponse) Scan(value interface{}) error {
	if value == nil {
		*r = LLMResponse{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, r)
}

// NewWorkPaperNote creates a new work paper note
func NewWorkPaperNote(workPaperID, masterItemID uuid.UUID) (*WorkPaperNote, error) {
	if workPaperID == uuid.Nil {
		return nil, ErrWorkPaperIDRequired
	}

	if masterItemID == uuid.Nil {
		return nil, ErrMasterItemIDRequired
	}

	now := time.Now()
	return &WorkPaperNote{
		ID:           uuid.New(),
		WorkPaperID:  workPaperID,
		MasterItemID: masterItemID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// UpdateGDriveLink updates the Google Drive link
func (wpn *WorkPaperNote) UpdateGDriveLink(link string) {
	if link == "" {
		wpn.GDriveLink = nil
	} else {
		wpn.GDriveLink = &link
	}
	wpn.UpdatedAt = time.Now()
}

// UpdateLLMResult updates the validation result from LLM
func (wpn *WorkPaperNote) UpdateLLMResult(isValid bool, notes string, response LLMResponse) {
	wpn.IsValid = &isValid
	if notes == "" {
		wpn.Notes = nil
	} else {
		wpn.Notes = &notes
	}
	if response.Note == "" && response.Model == "" && response.Usage == nil {
		wpn.LastLLMResponse = nil
	} else {
		wpn.LastLLMResponse = &response
	}
	wpn.UpdatedAt = time.Now()
}

// UpdateValidation allows manual override of validation result
func (wpn *WorkPaperNote) UpdateValidation(isValid *bool, notes string) {
	wpn.IsValid = isValid
	if notes == "" {
		wpn.Notes = nil
	} else {
		wpn.Notes = &notes
	}
	wpn.UpdatedAt = time.Now()
}

// Helper methods to safely access nullable fields

// GetGDriveLink safely returns the Google Drive link
func (wpn *WorkPaperNote) GetGDriveLink() string {
	if wpn.GDriveLink == nil {
		return ""
	}
	return *wpn.GDriveLink
}

// GetNotes safely returns the notes
func (wpn *WorkPaperNote) GetNotes() string {
	if wpn.Notes == nil {
		return ""
	}
	return *wpn.Notes
}

// GetLLMResponse safely returns the LLM response
func (wpn *WorkPaperNote) GetLLMResponse() LLMResponse {
	if wpn.LastLLMResponse == nil {
		return LLMResponse{}
	}
	return *wpn.LastLLMResponse
}