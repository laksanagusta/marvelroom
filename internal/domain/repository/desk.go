package repository

import (
	"context"

	"sandbox/internal/domain/entity"
	"sandbox/pkg/pagination"
)

// WorkPaperItemRepository defines the interface for work paper item data operations
type WorkPaperItemRepository interface {
	Create(ctx context.Context, item *entity.WorkPaperItem) (*entity.WorkPaperItem, error)
	GetByID(ctx context.Context, id string) (*entity.WorkPaperItem, error)
	GetByNumber(ctx context.Context, number string) (*entity.WorkPaperItem, error)
	Update(ctx context.Context, item *entity.WorkPaperItem) (*entity.WorkPaperItem, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, params *pagination.QueryParams) ([]*entity.WorkPaperItem, int64, error)
	ListActive(ctx context.Context) ([]*entity.WorkPaperItem, error)
}

// OrganizationRepository defines the interface for organization data operations
type OrganizationRepository interface {
	GetOrganizations(ctx context.Context, page, limit int, sort string) (*entity.OrganizationListResponse, error)
	GetByID(ctx context.Context, id string) (*entity.Organization, error)
}

// WorkPaperFilter defines the filter parameters for work paper queries
type WorkPaperFilter struct {
	Status         string `json:"status,omitempty"`
	OrganizationID string `json:"organization_id,omitempty"`
	Year           int    `json:"year,omitempty"`
	Semester       int    `json:"semester,omitempty"`
}

// WorkPaperRepository defines the interface for work paper data operations
type WorkPaperRepository interface {
	Create(ctx context.Context, wp *entity.WorkPaper) (*entity.WorkPaper, error)
	GetByID(ctx context.Context, id string) (*entity.WorkPaper, error)
	GetByOrganizationYearSemester(ctx context.Context, organizationID string, year, semester int) (*entity.WorkPaper, error)
	GetByFilter(ctx context.Context, filter *WorkPaperFilter, page, limit int) ([]*entity.WorkPaper, int64, error)
	Update(ctx context.Context, wp *entity.WorkPaper) (*entity.WorkPaper, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, params interface{}) ([]*entity.WorkPaper, int64, error)
	ListByOrganization(ctx context.Context, organizationID string) ([]*entity.WorkPaper, error)
}

// WorkPaperNoteRepository defines the interface for work paper note data operations
type WorkPaperNoteRepository interface {
	Create(ctx context.Context, note *entity.WorkPaperNote) (*entity.WorkPaperNote, error)
	CreateBatch(ctx context.Context, notes []*entity.WorkPaperNote) ([]*entity.WorkPaperNote, error)
	GetByID(ctx context.Context, id string) (*entity.WorkPaperNote, error)
	GetByWorkPaper(ctx context.Context, workPaperID string) ([]*entity.WorkPaperNote, error)
	Update(ctx context.Context, note *entity.WorkPaperNote) (*entity.WorkPaperNote, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, params interface{}) ([]*entity.WorkPaperNote, int64, error)
	WithTransaction(tx interface{}) WorkPaperNoteRepository
}

// Backward compatibility aliases (deprecated)
type MasterLakipItemRepository = WorkPaperItemRepository
type PaperWorkRepository = WorkPaperRepository
type PaperWorkItemRepository = WorkPaperNoteRepository