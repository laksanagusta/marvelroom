package entity

import (
	"time"

	"github.com/google/uuid"
)

// WorkPaperItem represents master data for work paper checklist items with tree structure
type WorkPaperItem struct {
	ID           uuid.UUID  `db:"id"`
	Type         string     `db:"type"`          // A, B, C type for work paper hierarchy
	Number       string     `db:"number"`        // Numbering like 1., 1.1, 1.1.1
	Statement    string     `db:"statement"`     // Pernyataan/Eksistensi
	Explanation  string     `db:"explanation"`   // Penjelasan
	FillingGuide string     `db:"filling_guide"` // Petunjuk Pengisian
	ParentID     *uuid.UUID `db:"parent_id"`     // Parent ID for tree structure
	Level        int        `db:"level"`         // Hierarchy level (1, 2, 3)
	SortOrder    int        `db:"sort_order"`    // Order within same level
	IsActive     bool       `db:"is_active"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`

	// Relations
	Children []*WorkPaperItem `db:"-"`
	Parent   *WorkPaperItem   `db:"-"`
}

// WorkPaperItemType constants
const (
	WorkPaperItemTypeA = "A" // Utama
	WorkPaperItemTypeB = "B" // Kegiatan Utama
	WorkPaperItemTypeC = "C" // Kegiatan Penunjang Utama
)

// NewWorkPaperItem creates a new work paper item with validation
func NewWorkPaperItem(itemType, number, statement, explanation, fillingGuide string, parentID *uuid.UUID, level, sortOrder int) (*WorkPaperItem, error) {
	if itemType == "" {
		return nil, ErrWorkPaperItemTypeRequired
	}

	if number == "" {
		return nil, ErrWorkPaperItemNumberRequired
	}

	if statement == "" {
		return nil, ErrWorkPaperItemStatementRequired
	}

	now := time.Now()
	return &WorkPaperItem{
		ID:           uuid.New(),
		Type:         itemType,
		Number:       number,
		Statement:    statement,
		Explanation:  explanation,
		FillingGuide: fillingGuide,
		ParentID:     parentID,
		Level:        level,
		SortOrder:    sortOrder,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Update updates the work paper item fields
func (w *WorkPaperItem) Update(itemType, number, statement, explanation, fillingGuide string, sortOrder *int) error {
	if itemType == "" {
		return ErrWorkPaperItemTypeRequired
	}

	if number == "" {
		return ErrWorkPaperItemNumberRequired
	}

	if statement == "" {
		return ErrWorkPaperItemStatementRequired
	}

	// Validate type
	if itemType != WorkPaperItemTypeA && itemType != WorkPaperItemTypeB && itemType != WorkPaperItemTypeC {
		return ErrInvalidWorkPaperItemType
	}

	w.Type = itemType
	w.Number = number
	w.Statement = statement
	w.Explanation = explanation
	w.FillingGuide = fillingGuide
	if sortOrder != nil {
		w.SortOrder = *sortOrder
	}
	w.UpdatedAt = time.Now()

	return nil
}

// Deactivate soft deletes the work paper item
func (w *WorkPaperItem) Deactivate() {
	now := time.Now()
	w.IsActive = false
	w.DeletedAt = &now
	w.UpdatedAt = now
}

// Activate reactivates a deactivated work paper item
func (w *WorkPaperItem) Activate() {
	w.IsActive = true
	w.DeletedAt = nil
	w.UpdatedAt = time.Now()
}

// IsRoot checks if this item is a root node (no parent)
func (w *WorkPaperItem) IsRoot() bool {
	return w.ParentID == nil
}

// IsLeaf checks if this item is a leaf node (no children)
func (w *WorkPaperItem) IsLeaf() bool {
	return len(w.Children) == 0
}

// AddChild adds a child item to this item
func (w *WorkPaperItem) AddChild(child *WorkPaperItem) {
	child.ParentID = &w.ID
	child.Level = w.Level + 1
	w.Children = append(w.Children, child)
}

// GetFullPath returns the full path numbering for this item
func (w *WorkPaperItem) GetFullPath() string {
	if w.Parent == nil {
		return w.Number
	}
	return w.Parent.GetFullPath() + "." + w.Number
}

// WorkPaperItemListResponse represents the response for listing work paper items
type WorkPaperItemListResponse struct {
	Data     []WorkPaperItem       `json:"data"`
	Metadata WorkPaperItemMetadata `json:"metadata"`
}

// WorkPaperItemMetadata represents pagination metadata for work paper items
type WorkPaperItemMetadata struct {
	Count       int `json:"count"`
	TotalCount  int `json:"total_count"`
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	PageSize    int `json:"page_size"`
}

// ListRequest represents the request parameters for listing work paper items
type ListRequest struct {
	Search   string `json:"search"`
	Type     string `json:"type"`
	IsActive *bool  `json:"is_active"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Sort     string `json:"sort"`
}

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page     int
	PageSize int
	Search   string
	Type     string
	IsActive *bool
	Sort     string
}
