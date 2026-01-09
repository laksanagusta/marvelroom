package entity

import (
	"time"

	"github.com/google/uuid"
)

// WorkPaperTopic represents a topic/classification for work paper items
type WorkPaperTopic struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	Name            string     `db:"name" json:"name"`
	Description     *string    `db:"description" json:"description,omitempty"`
	TemplatePath    *string    `db:"template_path" json:"template_path,omitempty"`
	TemplateVersion int        `db:"template_version" json:"template_version"`
	IsActive        bool       `db:"is_active" json:"is_active"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

// NewWorkPaperTopic creates a new work paper topic with validation
func NewWorkPaperTopic(name string, description, templatePath *string) (*WorkPaperTopic, error) {
	if name == "" {
		return nil, ErrWorkPaperTopicNameRequired
	}

	now := time.Now()
	return &WorkPaperTopic{
		ID:              uuid.New(),
		Name:            name,
		Description:     description,
		TemplatePath:    templatePath,
		TemplateVersion: 1,
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// Update updates the work paper topic fields
func (t *WorkPaperTopic) Update(name string, description, templatePath *string, isActive *bool) error {
	if name == "" {
		return ErrWorkPaperTopicNameRequired
	}

	t.Name = name
	t.Description = description
	t.TemplatePath = templatePath
	if isActive != nil {
		t.IsActive = *isActive
	}
	t.UpdatedAt = time.Now()

	return nil
}

// IncrementTemplateVersion increments the template version when template is updated
func (t *WorkPaperTopic) IncrementTemplateVersion() {
	t.TemplateVersion++
	t.UpdatedAt = time.Now()
}

// Deactivate soft deletes the work paper topic
func (t *WorkPaperTopic) Deactivate() {
	now := time.Now()
	t.IsActive = false
	t.DeletedAt = &now
	t.UpdatedAt = now
}

// Activate reactivates a deactivated work paper topic
func (t *WorkPaperTopic) Activate() {
	t.IsActive = true
	t.DeletedAt = nil
	t.UpdatedAt = time.Now()
}
