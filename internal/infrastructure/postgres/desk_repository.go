package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/pkg/database"
	"sandbox/pkg/pagination"
)

// Work paper item repository
type workPaperItemRepository struct {
	db database.Queryer
}

func NewWorkPaperItemRepository(db database.Queryer) repository.WorkPaperItemRepository {
	return &workPaperItemRepository{db: db}
}

func (r *workPaperItemRepository) Create(ctx context.Context, item *entity.WorkPaperItem) (*entity.WorkPaperItem, error) {
	query := `
		INSERT INTO work_paper_items (
			id, type, number, statement, explanation, filling_guide, parent_id, level, sort_order, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.ExecContext(ctx, query,
		item.ID, item.Type, item.Number, item.Statement, item.Explanation, item.FillingGuide,
		item.ParentID, item.Level, item.SortOrder, item.IsActive, item.CreatedAt, item.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create work paper item: %w", err)
	}
	return item, nil
}

func (r *workPaperItemRepository) GetByID(ctx context.Context, id string) (*entity.WorkPaperItem, error) {
	query := `
		SELECT id, type, number, statement, explanation, filling_guide, parent_id, level, sort_order, is_active, created_at, updated_at, deleted_at
		FROM work_paper_items
		WHERE id = $1 AND deleted_at IS NULL
	`

	var item entity.WorkPaperItem
	err := r.db.GetContext(ctx, &item, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrWorkPaperItemNotFound
		}
		return nil, fmt.Errorf("failed to get work paper item: %w", err)
	}
	return &item, nil
}

func (r *workPaperItemRepository) GetByNumber(ctx context.Context, number string) (*entity.WorkPaperItem, error) {
	query := `
		SELECT id, type, number, statement, explanation, filling_guide, parent_id, level, sort_order, is_active, created_at, updated_at, deleted_at
		FROM work_paper_items
		WHERE number = $1 AND deleted_at IS NULL
	`

	var item entity.WorkPaperItem
	err := r.db.GetContext(ctx, &item, query, number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrWorkPaperItemNotFound
		}
		return nil, fmt.Errorf("failed to get work paper item: %w", err)
	}
	return &item, nil
}

func (r *workPaperItemRepository) Update(ctx context.Context, item *entity.WorkPaperItem) (*entity.WorkPaperItem, error) {
	query := `
		UPDATE work_paper_items
		SET type = $2, number = $3, statement = $4, explanation = $5, filling_guide = $6, parent_id = $7, level = $8, sort_order = $9, updated_at = $10
		WHERE id = $1
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		item.ID, item.Type, item.Number, item.Statement, item.Explanation, item.FillingGuide,
		item.ParentID, item.Level, item.SortOrder, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update work paper item: %w", err)
	}
	item.UpdatedAt = now
	return item, nil
}

func (r *workPaperItemRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE work_paper_items
		SET deleted_at = $1, updated_at = $2, is_active = false
		WHERE id = $3
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, now, now, id)
	if err != nil {
		return fmt.Errorf("failed to delete work paper item: %w", err)
	}
	return nil
}

func (r *workPaperItemRepository) List(ctx context.Context, params *pagination.QueryParams) ([]*entity.WorkPaperItem, int64, error) {
	// Build count query
	countBuilder := pagination.NewQueryBuilder("SELECT COUNT(*) FROM work_paper_items")
	for _, filter := range params.Filters {
		if err := countBuilder.AddFilter(filter); err != nil {
			return nil, 0, err
		}
	}
	// Always include deleted_at filter
	countBuilder.AddFilter(pagination.Filter{
		Field:    "deleted_at",
		Operator: "is",
		Value:    nil,
	})
	countQuery, countArgs := countBuilder.Build()
	var totalCount int64
	err := r.db.GetContext(ctx, &totalCount, countQuery, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	// Build main query
	queryBuilder := pagination.NewQueryBuilder(`
		SELECT
			id, type, number, statement, explanation, filling_guide, parent_id, level, sort_order, is_active, created_at, updated_at, deleted_at
		FROM work_paper_items`)
	for _, filter := range params.Filters {
		if err := queryBuilder.AddFilter(filter); err != nil {
			return nil, 0, err
		}
	}
	// Always include deleted_at filter
	queryBuilder.AddFilter(pagination.Filter{
		Field:    "deleted_at",
		Operator: "is",
		Value:    nil,
	})

	// Add default sorting if no sorts provided
	if len(params.Sorts) == 0 {
		params.Sorts = []pagination.Sort{
			{Field: "level", Order: "asc"},
			{Field: "sort_order", Order: "asc"},
			{Field: "number", Order: "asc"},
		}
	}

	for _, sort := range params.Sorts {
		if err := queryBuilder.AddSort(sort); err != nil {
			return nil, 0, err
		}
	}
	query, args := queryBuilder.Build()

	// Add pagination
	offset := (params.Pagination.Page - 1) * params.Pagination.Limit
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", params.Pagination.Limit, offset)

	var workPaperItems []*entity.WorkPaperItem
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.WorkPaperItem
		if err := rows.StructScan(&item); err != nil {
			return nil, 0, fmt.Errorf("failed to scan work paper item: %w", err)
		}
		workPaperItems = append(workPaperItems, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return workPaperItems, totalCount, nil
}

func (r *workPaperItemRepository) ListActive(ctx context.Context) ([]*entity.WorkPaperItem, error) {
	query := `
		SELECT id, type, number, statement, explanation, filling_guide, parent_id, level, sort_order, is_active, created_at, updated_at, deleted_at
		FROM work_paper_items
		WHERE deleted_at IS NULL AND is_active = true
		ORDER BY level, sort_order, number ASC
	`

	var items []*entity.WorkPaperItem
	err := r.db.SelectContext(ctx, &items, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active work paper items: %w", err)
	}

	return items, nil
}

// Work paper repository
type workPaperRepository struct {
	db database.Queryer
}

func NewWorkPaperRepository(db database.Queryer) repository.WorkPaperRepository {
	return &workPaperRepository{db: db}
}

func (r *workPaperRepository) Create(ctx context.Context, wp *entity.WorkPaper) (*entity.WorkPaper, error) {
	log.Println(wp.OrganizationID)

	query := `
		INSERT INTO work_papers (
			id, organization_id, year, semester, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		wp.ID, wp.OrganizationID, wp.Year, wp.Semester, wp.Status, wp.CreatedAt, wp.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create work paper: %w", err)
	}
	return wp, nil
}

func (r *workPaperRepository) GetByID(ctx context.Context, id string) (*entity.WorkPaper, error) {
	query := `
		SELECT id, organization_id, year, semester, status, created_at, updated_at, deleted_at
		FROM work_papers
		WHERE id = $1 AND deleted_at IS NULL
	`

	var wp entity.WorkPaper
	err := r.db.GetContext(ctx, &wp, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrWorkPaperNotFound
		}
		return nil, fmt.Errorf("failed to get work paper: %w", err)
	}
	return &wp, nil
}

func (r *workPaperRepository) GetByOrganizationYearSemester(ctx context.Context, organizationID string, year, semester int) (*entity.WorkPaper, error) {
	query := `
		SELECT id, organization_id, year, semester, status, created_at, updated_at, deleted_at
		FROM work_papers
		WHERE organization_id = $1 AND year = $2 AND semester = $3 AND deleted_at IS NULL
	`

	var wp entity.WorkPaper
	err := r.db.GetContext(ctx, &wp, query, organizationID, year, semester)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrWorkPaperNotFound
		}
		return nil, fmt.Errorf("failed to get work paper: %w", err)
	}
	return &wp, nil
}

func (r *workPaperRepository) Update(ctx context.Context, wp *entity.WorkPaper) (*entity.WorkPaper, error) {
	query := `
		UPDATE work_papers
		SET status = $2, updated_at = $3
		WHERE id = $1
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, wp.ID, wp.Status, now)
	if err != nil {
		return nil, fmt.Errorf("failed to update work paper: %w", err)
	}
	wp.UpdatedAt = now
	return wp, nil
}

func (r *workPaperRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE work_papers
		SET deleted_at = $1, updated_at = $2
		WHERE id = $3
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, now, now, id)
	if err != nil {
		return fmt.Errorf("failed to delete work paper: %w", err)
	}
	return nil
}

func (r *workPaperRepository) List(ctx context.Context, params interface{}) ([]*entity.WorkPaper, int64, error) {
	query := `
		SELECT id, organization_id, year, semester, status, created_at, updated_at, deleted_at
		FROM work_papers
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	var wps []*entity.WorkPaper
	err := r.db.SelectContext(ctx, &wps, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query work papers: %w", err)
	}

	return wps, int64(len(wps)), nil
}

func (r *workPaperRepository) ListByOrganization(ctx context.Context, organizationID string) ([]*entity.WorkPaper, error) {
	query := `
		SELECT id, organization_id, year, semester, status, created_at, updated_at, deleted_at
		FROM work_papers
		WHERE organization_id = $1 AND deleted_at IS NULL
		ORDER BY year DESC, semester DESC
	`

	var wps []*entity.WorkPaper
	err := r.db.SelectContext(ctx, &wps, query, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to query work papers by organization: %w", err)
	}

	return wps, nil
}

func (r *workPaperRepository) GetByFilter(ctx context.Context, filter *repository.WorkPaperFilter, page, limit int) ([]*entity.WorkPaper, int64, error) {
	// Set default pagination values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	// Build count query
	countBuilder := pagination.NewQueryBuilder("SELECT COUNT(*) FROM work_papers")

	// Add deleted_at filter to count query
	countBuilder.AddFilter(pagination.Filter{
		Field:    "deleted_at",
		Operator: "is",
		Value:    nil,
	})

	// Add filter conditions to count query
	if filter != nil {
		if filter.Status != "" {
			countBuilder.AddFilter(pagination.Filter{
				Field:    "status",
				Operator: "eq",
				Value:    filter.Status,
			})
		}
		if filter.OrganizationID != "" {
			countBuilder.AddFilter(pagination.Filter{
				Field:    "organization_id",
				Operator: "eq",
				Value:    filter.OrganizationID,
			})
		}
		if filter.Year > 0 {
			countBuilder.AddFilter(pagination.Filter{
				Field:    "year",
				Operator: "eq",
				Value:    filter.Year,
			})
		}
		if filter.Semester > 0 {
			countBuilder.AddFilter(pagination.Filter{
				Field:    "semester",
				Operator: "eq",
				Value:    filter.Semester,
			})
		}
	}

	countQuery, countArgs := countBuilder.Build()
	var totalCount int64
	err := r.db.GetContext(ctx, &totalCount, countQuery, countArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count work papers: %w", err)
	}

	// Build main query
	queryBuilder := pagination.NewQueryBuilder(`
		SELECT
			id, organization_id, year, semester, status, created_at, updated_at, deleted_at
		FROM work_papers`)

	// Add deleted_at filter to main query
	queryBuilder.AddFilter(pagination.Filter{
		Field:    "deleted_at",
		Operator: "is",
		Value:    nil,
	})

	// Add filter conditions to main query
	if filter != nil {
		if filter.Status != "" {
			queryBuilder.AddFilter(pagination.Filter{
				Field:    "status",
				Operator: "eq",
				Value:    filter.Status,
			})
		}
		if filter.OrganizationID != "" {
			queryBuilder.AddFilter(pagination.Filter{
				Field:    "organization_id",
				Operator: "eq",
				Value:    filter.OrganizationID,
			})
		}
		if filter.Year > 0 {
			queryBuilder.AddFilter(pagination.Filter{
				Field:    "year",
				Operator: "eq",
				Value:    filter.Year,
			})
		}
		if filter.Semester > 0 {
			queryBuilder.AddFilter(pagination.Filter{
				Field:    "semester",
				Operator: "eq",
				Value:    filter.Semester,
			})
		}
	}

	// Add default ordering by created_at DESC
	queryBuilder.AddSort(pagination.Sort{
		Field: "created_at",
		Order: "desc",
	})

	query, args := queryBuilder.Build()

	// Add pagination
	offset := (page - 1) * limit
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	var workPapers []*entity.WorkPaper
	err = r.db.SelectContext(ctx, &workPapers, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query work papers with filter: %w", err)
	}

	return workPapers, totalCount, nil
}

// Work paper note repository
type workPaperNoteRepository struct {
	db database.Queryer
}

func NewWorkPaperNoteRepository(db database.Queryer) repository.WorkPaperNoteRepository {
	return &workPaperNoteRepository{db: db}
}

func (r *workPaperNoteRepository) Create(ctx context.Context, note *entity.WorkPaperNote) (*entity.WorkPaperNote, error) {
	query := `
		INSERT INTO work_paper_notes (
			id, work_paper_id, master_item_id, gdrive_link, is_valid, notes, last_llm_response, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		note.ID, note.WorkPaperID, note.MasterItemID, note.GDriveLink, note.IsValid,
		note.Notes, note.LastLLMResponse, note.CreatedAt, note.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create work paper note: %w", err)
	}
	return note, nil
}

func (r *workPaperNoteRepository) CreateBatch(ctx context.Context, notes []*entity.WorkPaperNote) ([]*entity.WorkPaperNote, error) {
	var createdNotes []*entity.WorkPaperNote
	for _, note := range notes {
		createdNote, err := r.Create(ctx, note)
		if err != nil {
			return nil, fmt.Errorf("failed to create work paper note: %w", err)
		}
		createdNotes = append(createdNotes, createdNote)
	}
	return createdNotes, nil
}

func (r *workPaperNoteRepository) GetByID(ctx context.Context, id string) (*entity.WorkPaperNote, error) {
	query := `
		SELECT id, work_paper_id, master_item_id, gdrive_link, is_valid, notes, last_llm_response, created_at, updated_at, deleted_at
		FROM work_paper_notes
		WHERE id = $1 AND deleted_at IS NULL
	`

	var note entity.WorkPaperNote
	err := r.db.GetContext(ctx, &note, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrWorkPaperNoteNotFound
		}
		return nil, fmt.Errorf("failed to get work paper note: %w", err)
	}

	// Load MasterItem data if available
	if note.MasterItemID != uuid.Nil {
		masterItem, err := r.getMasterItemByID(ctx, note.MasterItemID)
		if err != nil {
			log.Printf("Failed to load master item %s for note %s: %v", note.MasterItemID, note.ID, err)
			// Continue without master item if there's an error
			note.MasterItem = nil
		} else {
			note.MasterItem = masterItem
		}
	}

	return &note, nil
}

func (r *workPaperNoteRepository) GetByWorkPaper(ctx context.Context, workPaperID string) ([]*entity.WorkPaperNote, error) {
	query := `
		SELECT id, work_paper_id, master_item_id, gdrive_link, is_valid, notes, last_llm_response, created_at, updated_at, deleted_at
		FROM work_paper_notes
		WHERE work_paper_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC
	`

	var notes []*entity.WorkPaperNote
	err := r.db.SelectContext(ctx, &notes, query, workPaperID)
	if err != nil {
		return nil, fmt.Errorf("failed to query work paper notes: %w", err)
	}

	// Load MasterItem data for each note if available
	for _, note := range notes {
		if note.MasterItemID != uuid.Nil {
			masterItem, err := r.getMasterItemByID(ctx, note.MasterItemID)
			if err != nil {
				log.Printf("Failed to load master item %s for note %s: %v", note.MasterItemID, note.ID, err)
				// Continue without master item if there's an error
				note.MasterItem = nil
			} else {
				note.MasterItem = masterItem
			}
		}
	}

	return notes, nil
}

// getMasterItemByID is a helper method to load MasterItem data for a work paper note
func (r *workPaperNoteRepository) getMasterItemByID(ctx context.Context, masterItemID uuid.UUID) (*entity.WorkPaperItem, error) {
	query := `
		SELECT id, type, number, statement, explanation, filling_guide, parent_id, level, sort_order, is_active, created_at, updated_at, deleted_at
		FROM work_paper_items
		WHERE id = $1 AND deleted_at IS NULL
	`

	var masterItem entity.WorkPaperItem
	err := r.db.GetContext(ctx, &masterItem, query, masterItemID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrWorkPaperItemNotFound
		}
		return nil, fmt.Errorf("failed to get master item: %w", err)
	}
	return &masterItem, nil
}

func (r *workPaperNoteRepository) Update(ctx context.Context, note *entity.WorkPaperNote) (*entity.WorkPaperNote, error) {
	query := `
		UPDATE work_paper_notes
		SET gdrive_link = $2, is_valid = $3, notes = $4, last_llm_response = $5, updated_at = $6
		WHERE id = $1
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		note.ID, note.GDriveLink, note.IsValid, note.Notes, note.LastLLMResponse, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update work paper note: %w", err)
	}
	note.UpdatedAt = now
	return note, nil
}

func (r *workPaperNoteRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE work_paper_notes
		SET deleted_at = $1, updated_at = $2
		WHERE id = $3
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, now, now, id)
	if err != nil {
		return fmt.Errorf("failed to delete work paper note: %w", err)
	}
	return nil
}

func (r *workPaperNoteRepository) List(ctx context.Context, params interface{}) ([]*entity.WorkPaperNote, int64, error) {
	// Simplified implementation - can be expanded later
	return nil, 0, nil
}

func (r *workPaperNoteRepository) WithTransaction(tx interface{}) repository.WorkPaperNoteRepository {
	return r
}

// Backward compatibility factory functions (deprecated)
func NewMasterLakipItemRepository(db database.Queryer) repository.MasterLakipItemRepository {
	return NewWorkPaperItemRepository(db)
}

func NewPaperWorkRepository(db database.Queryer) repository.PaperWorkRepository {
	return NewWorkPaperRepository(db)
}

func NewPaperWorkItemRepository(db database.Queryer) repository.PaperWorkItemRepository {
	return NewWorkPaperNoteRepository(db)
}
