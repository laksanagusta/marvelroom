package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/pkg/pagination"
)

// deskService implements the DeskService interface
type deskService struct {
	workPaperItemRepo repository.WorkPaperItemRepository
	organizationRepo  repository.OrganizationRepository
	workPaperRepo     repository.WorkPaperRepository
	workPaperNoteRepo repository.WorkPaperNoteRepository
	signatureRepo     repository.WorkPaperSignatureRepository
	driveService      DriveService
	llmService        LLMService
}

// NewDeskService creates a new desk service instance
func NewDeskService(
	workPaperItemRepo repository.WorkPaperItemRepository,
	organizationRepo repository.OrganizationRepository,
	workPaperRepo repository.WorkPaperRepository,
	workPaperNoteRepo repository.WorkPaperNoteRepository,
	signatureRepo repository.WorkPaperSignatureRepository,
	driveService DriveService,
	llmService LLMService,
) DeskService {
	return &deskService{
		workPaperItemRepo: workPaperItemRepo,
		organizationRepo:  organizationRepo,
		workPaperRepo:     workPaperRepo,
		workPaperNoteRepo: workPaperNoteRepo,
		signatureRepo:     signatureRepo,
		driveService:      driveService,
		llmService:        llmService,
	}
}

// Work Paper Item operations

func (s *deskService) CreateWorkPaperItem(ctx context.Context, req *CreateWorkPaperItemRequest) (*entity.WorkPaperItem, error) {
	item, err := entity.NewWorkPaperItem(req.Type, req.Number, req.Statement, req.Explanation, req.FillingGuide, req.ParentID, req.Level, req.SortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to create work paper item: %w", err)
	}

	createdItem, err := s.workPaperItemRepo.Create(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("failed to save work paper item: %w", err)
	}

	return createdItem, nil
}

func (s *deskService) GetWorkPaperItem(ctx context.Context, id string) (*entity.WorkPaperItem, error) {
	item, err := s.workPaperItemRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get work paper item: %w", err)
	}

	return item, nil
}

func (s *deskService) UpdateWorkPaperItem(ctx context.Context, req *UpdateWorkPaperItemRequest) (*entity.WorkPaperItem, error) {
	item, err := s.workPaperItemRepo.GetByID(ctx, req.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get work paper item: %w", err)
	}

	// Update optional fields
	if req.IsActive != nil {
		if *req.IsActive {
			item.Activate()
		} else {
			item.Deactivate()
		}
	}

	// If IsActive is not provided, just update other fields
	err = item.Update(req.Type, req.Number, req.Statement, req.Explanation, req.FillingGuide, req.SortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to update work paper item: %w", err)
	}

	// Update parent ID if provided
	if req.ParentID != nil {
		item.ParentID = req.ParentID
		item.Level = req.Level
		// Also update the timestamp when parent changes
		item.UpdatedAt = time.Now()
	}

	updatedItem, err := s.workPaperItemRepo.Update(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("failed to save work paper item: %w", err)
	}

	return updatedItem, nil
}

func (s *deskService) DeleteWorkPaperItem(ctx context.Context, id uuid.UUID) (*entity.WorkPaperItem, error) {
	item, err := s.workPaperItemRepo.GetByID(ctx, id.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get work paper item: %w", err)
	}

	// Soft delete by deactivating
	item.Deactivate()

	updatedItem, err := s.workPaperItemRepo.Update(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("failed to delete work paper item: %w", err)
	}

	return updatedItem, nil
}

func (s *deskService) DeactivateWorkPaperItem(ctx context.Context, id string) error {
	item, err := s.workPaperItemRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get work paper item: %w", err)
	}

	item.Deactivate()

	_, err = s.workPaperItemRepo.Update(ctx, item)
	if err != nil {
		return fmt.Errorf("failed to deactivate work paper item: %w", err)
	}

	return nil
}

func (s *deskService) ActivateWorkPaperItem(ctx context.Context, id string) error {
	item, err := s.workPaperItemRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get work paper item: %w", err)
	}

	item.Activate()

	_, err = s.workPaperItemRepo.Update(ctx, item)
	if err != nil {
		return fmt.Errorf("failed to activate work paper item: %w", err)
	}

	return nil
}

func (s *deskService) ListWorkPaperItems(ctx context.Context, params *entity.PaginationParams) (*entity.WorkPaperItemListResponse, error) {
	// Convert old pagination params to new QueryParams format for backward compatibility
	queryParams := &pagination.QueryParams{
		Pagination: pagination.Pagination{
			Page:  params.Page,
			Limit: params.PageSize,
		},
		Filters: []pagination.Filter{},
		Sorts:   []pagination.Sort{},
	}

	// Note: Search functionality with multiple fields needs custom query building
	// For now, we'll search by statement only when search is provided
	if params.Search != "" {
		queryParams.Filters = append(queryParams.Filters, pagination.Filter{
			Field:    "statement",
			Operator: "ilike",
			Value:    params.Search,
		})
	}

	// Add type filter if provided
	if params.Type != "" {
		queryParams.Filters = append(queryParams.Filters, pagination.Filter{
			Field:    "type",
			Operator: "eq",
			Value:    params.Type,
		})
	}

	// Add is_active filter if provided
	if params.IsActive != nil {
		queryParams.Filters = append(queryParams.Filters, pagination.Filter{
			Field:    "is_active",
			Operator: "eq",
			Value:    *params.IsActive,
		})
	}

	// Parse sort if provided
	if params.Sort != "" {
		// Simple sort parsing - in a real implementation you might want more sophisticated parsing
		sorts := strings.Split(params.Sort, ",")
		for _, sort := range sorts {
			sort = strings.TrimSpace(sort)
			if strings.HasSuffix(sort, " DESC") {
				field := strings.TrimSuffix(sort, " DESC")
				queryParams.Sorts = append(queryParams.Sorts, pagination.Sort{
					Field: field,
					Order: "desc",
				})
			} else {
				field := strings.TrimSuffix(sort, " ASC")
				queryParams.Sorts = append(queryParams.Sorts, pagination.Sort{
					Field: field,
					Order: "asc",
				})
			}
		}
	}

	workPaperItems, totalCount, err := s.workPaperItemRepo.List(ctx, queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to list work paper items: %w", err)
	}

	// Convert to response format
	responseItems := make([]entity.WorkPaperItem, len(workPaperItems))
	for i, item := range workPaperItems {
		responseItems[i] = *item
	}

	// Calculate pagination metadata
	totalPages := (int(totalCount) + params.PageSize - 1) / params.PageSize
	metadata := entity.WorkPaperItemMetadata{
		Count:       len(responseItems),
		TotalCount:  int(totalCount),
		CurrentPage: params.Page,
		TotalPage:   totalPages,
		PageSize:    params.PageSize,
	}

	return &entity.WorkPaperItemListResponse{
		Data:     responseItems,
		Metadata: metadata,
	}, nil
}

func (s *deskService) GetActiveWorkPaperItems(ctx context.Context) ([]*entity.WorkPaperItem, error) {
	items, err := s.workPaperItemRepo.ListActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active work paper items: %w", err)
	}

	return items, nil
}

// Organization operations

func (s *deskService) GetOrganizations(ctx context.Context, page, limit int, sort string) (*entity.OrganizationListResponse, error) {
	orgs, err := s.organizationRepo.GetOrganizations(ctx, page, limit, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizations: %w", err)
	}

	return orgs, nil
}

func (s *deskService) GetOrganization(ctx context.Context, id string) (*entity.Organization, error) {
	org, err := s.organizationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	return org, nil
}

// Work Paper operations

func (s *deskService) CreateWorkPaper(ctx context.Context, req *CreateWorkPaperRequest) (*entity.WorkPaper, error) {
	organizationID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	// Check if organization exists
	_, err = s.organizationRepo.GetByID(ctx, req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}

	// Check if work paper already exists for this organization, year, and semester
	existingWorkPaper, _ := s.workPaperRepo.GetByOrganizationYearSemester(ctx, req.OrganizationID, req.Year, req.Semester)
	if existingWorkPaper != nil {
		return nil, entity.ErrDuplicateWorkPaper
	}

	workPaper, err := entity.NewWorkPaper(organizationID, req.Year, req.Semester)
	if err != nil {
		return nil, fmt.Errorf("failed to create work paper: %w", err)
	}

	createdWorkPaper, err := s.workPaperRepo.Create(ctx, workPaper)
	if err != nil {
		return nil, fmt.Errorf("failed to save work paper: %w", err)
	}

	// Create work paper notes from active master items
	masterItems, err := s.workPaperItemRepo.ListActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get master items: %w", err)
	}

	log.Printf("Found %d active work paper items, creating notes for work paper %s", len(masterItems), createdWorkPaper.ID.String())

	var notes []*entity.WorkPaperNote
	for _, masterItem := range masterItems {
		note, err := entity.NewWorkPaperNote(createdWorkPaper.ID, masterItem.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to create work paper note: %w", err)
		}
		notes = append(notes, note)
	}

	// Batch create notes
	if len(notes) > 0 {
		createdNotes, err := s.workPaperNoteRepo.CreateBatch(ctx, notes)
		if err != nil {
			return nil, fmt.Errorf("failed to create work paper notes: %w", err)
		}
		log.Printf("Successfully created %d work paper notes for work paper %s", len(createdNotes), createdWorkPaper.ID.String())
	} else {
		log.Printf("No active work paper items found, no notes created for work paper %s", createdWorkPaper.ID.String())
	}

	return createdWorkPaper, nil
}

func (s *deskService) GetWorkPaper(ctx context.Context, id string) (*entity.WorkPaper, error) {
	workPaper, err := s.workPaperRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get work paper: %w", err)
	}

	return workPaper, nil
}

func (s *deskService) GetWorkPaperByOrganizationYearSemester(ctx context.Context, organizationID string, year, semester int) (*entity.WorkPaper, error) {
	workPaper, err := s.workPaperRepo.GetByOrganizationYearSemester(ctx, organizationID, year, semester)
	if err != nil {
		return nil, fmt.Errorf("failed to get work paper: %w", err)
	}

	return workPaper, nil
}

func (s *deskService) UpdateWorkPaperStatus(ctx context.Context, id string, status string) error {
	workPaper, err := s.workPaperRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get work paper: %w", err)
	}

	err = workPaper.UpdateStatus(status)
	if err != nil {
		return fmt.Errorf("failed to update work paper status: %w", err)
	}

	_, err = s.workPaperRepo.Update(ctx, workPaper)
	if err != nil {
		return fmt.Errorf("failed to save work paper: %w", err)
	}

	return nil
}

func (s *deskService) ListWorkPapers(ctx context.Context, req *ListWorkPapersRequest) ([]*entity.WorkPaper, int64, error) {
	// Simplified implementation for now
	workPapers, total, err := s.workPaperRepo.List(ctx, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list work papers: %w", err)
	}

	// Apply filters
	if req.OrganizationID != "" || req.Year != nil || req.Semester != nil || req.Status != "" {
		var filteredWorkPapers []*entity.WorkPaper
		for _, workPaper := range workPapers {
			if req.OrganizationID != "" && workPaper.OrganizationID.String() != req.OrganizationID {
				continue
			}
			if req.Year != nil && workPaper.Year != *req.Year {
				continue
			}
			if req.Semester != nil && workPaper.Semester != *req.Semester {
				continue
			}
			if req.Status != "" && workPaper.Status != req.Status {
				continue
			}
			filteredWorkPapers = append(filteredWorkPapers, workPaper)
		}
		return filteredWorkPapers, int64(len(filteredWorkPapers)), nil
	}

	return workPapers, total, nil
}

func (s *deskService) ListWorkPapersByOrganization(ctx context.Context, organizationID string) ([]*entity.WorkPaper, error) {
	workPapers, err := s.workPaperRepo.ListByOrganization(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to list work papers by organization: %w", err)
	}

	return workPapers, nil
}

// Work Paper Note operations

func (s *deskService) GetWorkPaperNotes(ctx context.Context, workPaperID string) ([]*entity.WorkPaperNote, error) {
	notes, err := s.workPaperNoteRepo.GetByWorkPaper(ctx, workPaperID)
	if err != nil {
		return nil, fmt.Errorf("failed to get work paper notes: %w", err)
	}

	return notes, nil
}

func (s *deskService) GetWorkPaperNoteByID(ctx context.Context, noteID string) (*entity.WorkPaperNote, error) {
	note, err := s.workPaperNoteRepo.GetByID(ctx, noteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get work paper note: %w", err)
	}
	return note, nil
}

func (s *deskService) UpdateWorkPaperNoteLink(ctx context.Context, noteID string, driveLink string) (*entity.WorkPaperNote, error) {
	note, err := s.workPaperNoteRepo.GetByID(ctx, noteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get work paper note: %w", err)
	}

	note.UpdateGDriveLink(driveLink)

	updatedNote, err := s.workPaperNoteRepo.Update(ctx, note)
	if err != nil {
		return nil, fmt.Errorf("failed to update work paper note: %w", err)
	}

	return updatedNote, nil
}

func (s *deskService) CheckDocument(ctx context.Context, noteID string) (*CheckDocumentResponse, error) {
	note, err := s.workPaperNoteRepo.GetByID(ctx, noteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get work paper note: %w", err)
	}

	masterItem, err := s.workPaperItemRepo.GetByID(ctx, note.MasterItemID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get master item: %w", err)
	}

	var documents []DocumentFile
	if note.GetGDriveLink() != "" {
		gdriveLink := note.GetGDriveLink()
		log.Printf("Processing Google Drive link: %s", gdriveLink)
		files, err := s.driveService.GetFilesFromFolder(ctx, gdriveLink)
		if err != nil {
			return nil, fmt.Errorf("failed to get files from Google Drive: %w", err)
		}

		log.Printf("Found %d files from Google Drive", len(files))
		for _, file := range files {
			log.Printf("Downloading file: %s (ID: %s, Type: %s)", file.Name, file.ID, file.Type)
			data, err := s.driveService.DownloadFile(ctx, file.ID)
			if err != nil {
				log.Printf("Failed to download file %s: %v", file.Name, err)
				// Log error but continue with other files
				continue
			}
			log.Printf("Successfully downloaded file %s (%d bytes)", file.Name, len(data))
			documents = append(documents, DocumentFile{
				Name: file.Name,
				Data: data,
				Type: file.Type,
			})
		}
	} else {
		log.Printf("No Google Drive link found for note ID: %s", noteID)
	}

	log.Printf("Total documents prepared for LLM: %d", len(documents))

	// Check documents using LLM
	llmReq := &DocumentCheckRequest{
		Number:       masterItem.Number,
		Statement:    masterItem.Statement,
		Explanation:  masterItem.Explanation,
		FillingGuide: masterItem.FillingGuide,
		Documents:    documents,
	}

	llmResp, err := s.llmService.CheckDocument(ctx, llmReq)
	if err != nil {
		return nil, fmt.Errorf("failed to check document: %w", err)
	}

	// Update work paper note with LLM response
	var usage *entity.Usage
	if llmResp.Usage != nil {
		usage = &entity.Usage{
			PromptTokens:     llmResp.Usage.PromptTokens,
			CompletionTokens: llmResp.Usage.CompletionTokens,
			TotalTokens:      llmResp.Usage.TotalTokens,
		}
	}

	llmResponseData := entity.LLMResponse{
		Note:    llmResp.Notes,
		IsValid: llmResp.IsValid,
		Model:   llmResp.Model,
		Usage:   usage,
	}

	note.UpdateLLMResult(llmResp.IsValid, llmResp.Notes, llmResponseData)

	_, err = s.workPaperNoteRepo.Update(ctx, note)
	if err != nil {
		return nil, fmt.Errorf("failed to update work paper note: %w", err)
	}

	return &CheckDocumentResponse{
		IsValid: llmResp.IsValid,
		Notes:   llmResp.Notes,
		Model:   llmResp.Model,
	}, nil
}

func (s *deskService) UpdateWorkPaperNoteValidation(ctx context.Context, noteID string, isValid *bool, notes string) (*entity.WorkPaperNote, error) {
	note, err := s.workPaperNoteRepo.GetByID(ctx, noteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get work paper note: %w", err)
	}

	note.UpdateValidation(isValid, notes)

	updatedNote, err := s.workPaperNoteRepo.Update(ctx, note)
	if err != nil {
		return nil, fmt.Errorf("failed to update work paper note: %w", err)
	}

	return updatedNote, nil
}

// Backward compatibility methods (deprecated)
// Note: These methods are provided for backward compatibility but will be removed in future versions
// Please use the new WorkPaper* methods instead

func (s *deskService) CreateMasterLakipItem(ctx context.Context, req *CreateMasterLakipItemRequest) (*entity.WorkPaperItem, error) {
	// Convert to new request type and call new method
	newReq := &CreateWorkPaperItemRequest{
		Type:         req.Type,
		Number:       req.Number,
		Statement:    req.Statement,
		Explanation:  req.Explanation,
		FillingGuide: req.FillingGuide,
		ParentID:     req.ParentID,
		Level:        req.Level,
		SortOrder:    req.SortOrder,
	}
	return s.CreateWorkPaperItem(ctx, newReq)
}

func (s *deskService) GetMasterLakipItem(ctx context.Context, id string) (*entity.WorkPaperItem, error) {
	return s.GetWorkPaperItem(ctx, id)
}

func (s *deskService) UpdateMasterLakipItem(ctx context.Context, id string, req *UpdateMasterLakipItemRequest) (*entity.WorkPaperItem, error) {
	// Convert to new request type and call new method
	newReq := &UpdateWorkPaperItemRequest{
		Type:         req.Type,
		Number:       req.Number,
		Statement:    req.Statement,
		Explanation:  req.Explanation,
		FillingGuide: req.FillingGuide,
		SortOrder:    req.SortOrder,
	}
	// Set ID from parameter
	newReq.ID, _ = uuid.Parse(id)
	return s.UpdateWorkPaperItem(ctx, newReq)
}

func (s *deskService) DeactivateMasterLakipItem(ctx context.Context, id string) error {
	return s.DeactivateWorkPaperItem(ctx, id)
}

func (s *deskService) ActivateMasterLakipItem(ctx context.Context, id string) error {
	return s.ActivateWorkPaperItem(ctx, id)
}

func (s *deskService) ListMasterLakipItems(ctx context.Context, params *entity.PaginationParams) (*entity.WorkPaperItemListResponse, error) {
	// Call the new method directly for backward compatibility
	return s.ListWorkPaperItems(ctx, params)
}

func (s *deskService) GetActiveMasterLakipItems(ctx context.Context) ([]*entity.WorkPaperItem, error) {
	return s.GetActiveWorkPaperItems(ctx)
}

func (s *deskService) CreatePaperWork(ctx context.Context, req *CreatePaperWorkRequest) (*entity.WorkPaper, error) {
	// Convert to new request type and call new method
	newReq := &CreateWorkPaperRequest{
		OrganizationID: req.OrganizationID,
		Year:           req.Year,
		Semester:       req.Semester,
	}
	return s.CreateWorkPaper(ctx, newReq)
}

func (s *deskService) GetPaperWork(ctx context.Context, id string) (*entity.WorkPaper, error) {
	return s.workPaperRepo.GetByID(ctx, id)
}

func (s *deskService) GetPaperWorkByOrganizationYearSemester(ctx context.Context, organizationID string, year, semester int) (*entity.WorkPaper, error) {
	return s.workPaperRepo.GetByOrganizationYearSemester(ctx, organizationID, year, semester)
}

func (s *deskService) UpdatePaperWorkStatus(ctx context.Context, id string, status string) error {
	workPaper, err := s.workPaperRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get work paper: %w", err)
	}

	err = workPaper.UpdateStatus(status)
	if err != nil {
		return fmt.Errorf("failed to update work paper status: %w", err)
	}

	_, err = s.workPaperRepo.Update(ctx, workPaper)
	if err != nil {
		return fmt.Errorf("failed to save work paper: %w", err)
	}

	return nil
}

func (s *deskService) ListPaperWorks(ctx context.Context, req *ListPaperWorksRequest) ([]*entity.WorkPaper, int64, error) {
	// Simplified implementation for now
	workPapers, total, err := s.workPaperRepo.List(ctx, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list work papers: %w", err)
	}

	// Apply filters
	if req.OrganizationID != "" || req.Year != nil || req.Semester != nil || req.Status != "" {
		var filteredWorkPapers []*entity.WorkPaper
		for _, workPaper := range workPapers {
			if req.OrganizationID != "" && workPaper.OrganizationID.String() != req.OrganizationID {
				continue
			}
			if req.Year != nil && workPaper.Year != *req.Year {
				continue
			}
			if req.Semester != nil && workPaper.Semester != *req.Semester {
				continue
			}
			if req.Status != "" && workPaper.Status != req.Status {
				continue
			}
			filteredWorkPapers = append(filteredWorkPapers, workPaper)
		}
		return filteredWorkPapers, int64(len(filteredWorkPapers)), nil
	}

	return workPapers, total, nil
}

func (s *deskService) ListPaperWorksByOrganization(ctx context.Context, organizationID string) ([]*entity.WorkPaper, error) {
	return s.workPaperRepo.ListByOrganization(ctx, organizationID)
}

func (s *deskService) UpdatePaperWorkItemLink(ctx context.Context, itemID string, driveLink string) (*entity.WorkPaperNote, error) {
	note, err := s.workPaperNoteRepo.GetByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get work paper note: %w", err)
	}

	note.UpdateGDriveLink(driveLink)

	updatedNote, err := s.workPaperNoteRepo.Update(ctx, note)
	if err != nil {
		return nil, fmt.Errorf("failed to update work paper note: %w", err)
	}

	return updatedNote, nil
}

func (s *deskService) UpdatePaperWorkItemValidation(ctx context.Context, itemID string, isValid *bool, notes string) (*entity.WorkPaperNote, error) {
	note, err := s.workPaperNoteRepo.GetByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get work paper note: %w", err)
	}

	note.UpdateValidation(isValid, notes)

	updatedNote, err := s.workPaperNoteRepo.Update(ctx, note)
	if err != nil {
		return nil, fmt.Errorf("failed to update work paper note: %w", err)
	}

	return updatedNote, nil
}

// Work Paper Signature operations

func (s *deskService) CreateWorkPaperSignature(ctx context.Context, req *CreateWorkPaperSignatureRequest) (*entity.WorkPaperSignature, error) {
	// Parse work paper ID
	workPaperID, err := uuid.Parse(req.WorkPaperID)
	if err != nil {
		return nil, fmt.Errorf("invalid work paper ID: %w", err)
	}

	// Check if work paper exists
	_, err = s.workPaperRepo.GetByID(ctx, req.WorkPaperID)
	if err != nil {
		return nil, fmt.Errorf("work paper not found: %w", err)
	}

	// Check if signature already exists for this user and paper
	_, err = s.signatureRepo.GetByWorkPaperIDAndUserID(ctx, workPaperID, req.UserID)
	if err == nil {
		return nil, entity.ErrDuplicateSignature
	} else if err != entity.ErrSignatureNotFound {
		return nil, fmt.Errorf("failed to check existing signature: %w", err)
	}

	// Create new signature
	signature, err := entity.NewWorkPaperSignature(workPaperID, req.UserID, req.UserName, req.SignatureType)
	if err != nil {
		return nil, fmt.Errorf("failed to create signature: %w", err)
	}

	// Set user details
	signature.SetUserDetails(req.UserEmail, req.UserRole)

	// Add signature data if provided
	if req.SignatureData != nil {
		signature.AddSignatureData(*req.SignatureData)
	}

	// Save to database
	err = s.signatureRepo.Create(ctx, signature)
	if err != nil {
		return nil, fmt.Errorf("failed to create signature in database: %w", err)
	}

	log.Printf("Created new work paper signature: ID=%s, WorkPaperID=%s, UserID=%s",
		signature.ID, signature.WorkPaperID, signature.UserID)

	return signature, nil
}

func (s *deskService) GetWorkPaperSignature(ctx context.Context, signatureID string) (*entity.WorkPaperSignature, error) {
	// Parse signature ID
	id, err := uuid.Parse(signatureID)
	if err != nil {
		return nil, fmt.Errorf("invalid signature ID: %w", err)
	}

	// Get signature from database
	signature, err := s.signatureRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get signature: %w", err)
	}

	return signature, nil
}

func (s *deskService) GetWorkPaperSignatures(ctx context.Context, paperID string) ([]*entity.WorkPaperSignature, error) {
	// Parse work paper ID
	workPaperID, err := uuid.Parse(paperID)
	if err != nil {
		return nil, fmt.Errorf("invalid work paper ID: %w", err)
	}

	// Get all signatures for the paper
	signatures, err := s.signatureRepo.GetByWorkPaperID(ctx, workPaperID)
	if err != nil {
		return nil, fmt.Errorf("failed to get signatures: %w", err)
	}

	return signatures, nil
}

func (s *deskService) SignWorkPaper(ctx context.Context, signatureID string, req *SignWorkPaperRequest) (*entity.WorkPaperSignature, error) {
	// Parse signature ID
	id, err := uuid.Parse(signatureID)
	if err != nil {
		return nil, fmt.Errorf("invalid signature ID: %w", err)
	}

	// Get signature from database
	signature, err := s.signatureRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get signature: %w", err)
	}

	// Add signature data if provided
	if req.SignatureData != nil {
		signature.AddSignatureData(*req.SignatureData)
	}

	// Sign the work paper
	err = signature.Sign(req.Notes)
	if err != nil {
		return nil, fmt.Errorf("failed to sign work paper: %w", err)
	}

	// Update in database
	err = s.signatureRepo.Update(ctx, signature)
	if err != nil {
		return nil, fmt.Errorf("failed to update signature: %w", err)
	}

	log.Printf("Work paper signed: SignatureID=%s, UserID=%s, SignedAt=%v",
		signature.ID, signature.UserID, signature.GetSignedAt())

	return signature, nil
}

func (s *deskService) RejectWorkPaperSignature(ctx context.Context, signatureID string, req *RejectWorkPaperSignatureRequest) (*entity.WorkPaperSignature, error) {
	// Parse signature ID
	id, err := uuid.Parse(signatureID)
	if err != nil {
		return nil, fmt.Errorf("invalid signature ID: %w", err)
	}

	// Get signature from database
	signature, err := s.signatureRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get signature: %w", err)
	}

	// Reject the signature
	err = signature.Reject(req.Notes)
	if err != nil {
		return nil, fmt.Errorf("failed to reject signature: %w", err)
	}

	// Update in database
	err = s.signatureRepo.Update(ctx, signature)
	if err != nil {
		return nil, fmt.Errorf("failed to update signature: %w", err)
	}

	log.Printf("Work paper signature rejected: SignatureID=%s, UserID=%s",
		signature.ID, signature.UserID)

	return signature, nil
}

func (s *deskService) ResetWorkPaperSignature(ctx context.Context, signatureID string) (*entity.WorkPaperSignature, error) {
	// Parse signature ID
	id, err := uuid.Parse(signatureID)
	if err != nil {
		return nil, fmt.Errorf("invalid signature ID: %w", err)
	}

	// Get signature from database
	signature, err := s.signatureRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get signature: %w", err)
	}

	// Reset the signature
	err = signature.Reset()
	if err != nil {
		return nil, fmt.Errorf("failed to reset signature: %w", err)
	}

	// Update in database
	err = s.signatureRepo.Update(ctx, signature)
	if err != nil {
		return nil, fmt.Errorf("failed to update signature: %w", err)
	}

	log.Printf("Work paper signature reset: SignatureID=%s, UserID=%s",
		signature.ID, signature.UserID)

	return signature, nil
}

func (s *deskService) GetWorkPaperSignaturesByUserID(ctx context.Context, userID string) ([]*entity.WorkPaperSignature, error) {
	// Get all signatures by user ID
	signatures, err := s.signatureRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get signatures by user ID: %w", err)
	}

	return signatures, nil
}

func (s *deskService) GetPendingSignaturesByUserID(ctx context.Context, userID string) ([]*entity.WorkPaperSignature, error) {
	// Get pending signatures by user ID
	signatures, err := s.signatureRepo.GetPendingByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending signatures by user ID: %w", err)
	}

	return signatures, nil
}

func (s *deskService) GetPendingSignaturesByPaperID(ctx context.Context, paperID string) ([]*entity.WorkPaperSignature, error) {
	// Parse work paper ID
	workPaperID, err := uuid.Parse(paperID)
	if err != nil {
		return nil, fmt.Errorf("invalid work paper ID: %w", err)
	}

	// Get pending signatures for the paper
	signatures, err := s.signatureRepo.GetPendingSignatures(ctx, workPaperID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending signatures: %w", err)
	}

	return signatures, nil
}

func (s *deskService) GetSignatureStatsByNoteID(ctx context.Context, paperID string) (*SignatureStatsResponse, error) {
	// Parse work paper ID
	workPaperID, err := uuid.Parse(paperID)
	if err != nil {
		return nil, fmt.Errorf("invalid work paper ID: %w", err)
	}

	// Get signature statistics
	stats, err := s.signatureRepo.GetSignatureStats(ctx, workPaperID)
	if err != nil {
		return nil, fmt.Errorf("failed to get signature statistics: %w", err)
	}

	return &SignatureStatsResponse{
		Total:    stats.Total,
		Pending:  stats.Pending,
		Signed:   stats.Signed,
		Rejected: stats.Rejected,
	}, nil
}

func (s *deskService) GetWorkPapersWithSignatures(ctx context.Context, page, limit int, status, organizationID string) ([]*WorkPaperWithSignatures, error) {
	// Create filter for work papers
	filter := &repository.WorkPaperFilter{
		Status:         status,
		OrganizationID: organizationID,
	}

	// Get work papers with pagination
	workPapers, _, err := s.workPaperRepo.GetByFilter(ctx, filter, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get work papers: %w", err)
	}

	// Build response with signatures
	var workPapersWithSignatures []*WorkPaperWithSignatures
	for _, workPaper := range workPapers {
		// Get signatures for this work paper
		signatures, err := s.signatureRepo.GetByWorkPaperID(ctx, workPaper.ID)
		if err != nil {
			log.Printf("Failed to get signatures for work paper %s: %v", workPaper.ID.String(), err)
			// Continue with empty signatures if there's an error
			signatures = []*entity.WorkPaperSignature{}
		}

		// Create WorkPaperWithSignatures
		workPaperWithSignatures := &WorkPaperWithSignatures{
			WorkPaper:  workPaper,
			Signatures: signatures,
		}

		workPapersWithSignatures = append(workPapersWithSignatures, workPaperWithSignatures)
	}

	log.Printf("Retrieved %d work papers with signatures", len(workPapersWithSignatures))
	return workPapersWithSignatures, nil
}

func (s *deskService) ManageSigners(ctx context.Context, req *ManageSignersRequest) (*ManageSignersResponse, error) {
	// Parse work paper ID
	workPaperID, err := uuid.Parse(req.WorkPaperID)
	if err != nil {
		return nil, fmt.Errorf("invalid work paper ID: %w", err)
	}

	// Check if work paper exists
	_, err = s.workPaperRepo.GetByID(ctx, req.WorkPaperID)
	if err != nil {
		return nil, fmt.Errorf("work paper not found: %w", err)
	}

	var signerResponses []SignerResponse

	switch req.Action {
	case "add":
		for _, signerData := range req.Signers {
			// Check if signature already exists
			_, err := s.signatureRepo.GetByWorkPaperIDAndUserID(ctx, workPaperID, signerData.UserID)
			if err == nil {
				// Signature already exists, skip
				continue
			} else if err != entity.ErrSignatureNotFound {
				return nil, fmt.Errorf("failed to check existing signature: %w", err)
			}

			// Create new signature
			signature, err := entity.NewWorkPaperSignature(workPaperID, signerData.UserID, signerData.UserName, signerData.SignatureType)
			if err != nil {
				return nil, fmt.Errorf("failed to create signature: %w", err)
			}

			// Set user details
			signature.SetUserDetails(signerData.UserEmail, signerData.UserRole)

			// Save to database
			err = s.signatureRepo.Create(ctx, signature)
			if err != nil {
				return nil, fmt.Errorf("failed to create signature: %w", err)
			}

			signerResponses = append(signerResponses, SignerResponse{
				SignatureID:   signature.ID.String(),
				UserID:        signature.UserID,
				UserName:      signature.UserName,
				UserEmail:     signature.GetUserEmail(),
				UserRole:      signature.GetUserRole(),
				SignatureType: signature.SignatureType,
				Status:        signature.Status,
				CreatedAt:     signature.CreatedAt.Format("2006-01-02T15:04:05Z"),
			})
		}

	case "remove":
		for _, signerData := range req.Signers {
			// Get existing signature
			signature, err := s.signatureRepo.GetByWorkPaperIDAndUserID(ctx, workPaperID, signerData.UserID)
			if err != nil {
				if err == entity.ErrSignatureNotFound {
					continue // Skip if not found
				}
				return nil, fmt.Errorf("failed to get signature: %w", err)
			}

			// Delete signature
			err = s.signatureRepo.Delete(ctx, signature.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to delete signature: %w", err)
			}
		}

	case "replace":
		// First, remove all existing signatures for the work paper
		existingSignatures, err := s.signatureRepo.GetByWorkPaperID(ctx, workPaperID)
		if err != nil {
			return nil, fmt.Errorf("failed to get existing signatures: %w", err)
		}

		for _, existingSig := range existingSignatures {
			err = s.signatureRepo.Delete(ctx, existingSig.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to delete existing signature: %w", err)
			}
		}

		// Then add new signatures
		for _, signerData := range req.Signers {
			signature, err := entity.NewWorkPaperSignature(workPaperID, signerData.UserID, signerData.UserName, signerData.SignatureType)
			if err != nil {
				return nil, fmt.Errorf("failed to create signature: %w", err)
			}

			signature.SetUserDetails(signerData.UserEmail, signerData.UserRole)

			err = s.signatureRepo.Create(ctx, signature)
			if err != nil {
				return nil, fmt.Errorf("failed to create signature: %w", err)
			}

			signerResponses = append(signerResponses, SignerResponse{
				SignatureID:   signature.ID.String(),
				UserID:        signature.UserID,
				UserName:      signature.UserName,
				UserEmail:     signature.GetUserEmail(),
				UserRole:      signature.GetUserRole(),
				SignatureType: signature.SignatureType,
				Status:        signature.Status,
				CreatedAt:     signature.CreatedAt.Format("2006-01-02T15:04:05Z"),
			})
		}

	default:
		return nil, fmt.Errorf("invalid action: %s", req.Action)
	}

	return &ManageSignersResponse{
		WorkPaperID: req.WorkPaperID,
		Action:      req.Action,
		Signers:     signerResponses,
		Message:     fmt.Sprintf("Successfully completed %s action on signers", req.Action),
	}, nil
}

func (s *deskService) SignWorkPaperWithUser(ctx context.Context, signatureID string, userID string) (*entity.WorkPaperSignature, error) {
	// Parse signature ID
	id, err := uuid.Parse(signatureID)
	if err != nil {
		return nil, fmt.Errorf("invalid signature ID: %w", err)
	}

	// Get signature from database
	signature, err := s.signatureRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get signature: %w", err)
	}

	// Verify that the signature belongs to the specified user
	if signature.UserID != userID {
		return nil, fmt.Errorf("signature does not belong to user: %s", userID)
	}

	// Sign the work paper
	err = signature.Sign("")
	if err != nil {
		return nil, fmt.Errorf("failed to sign work paper: %w", err)
	}

	// Update in database
	err = s.signatureRepo.Update(ctx, signature)
	if err != nil {
		return nil, fmt.Errorf("failed to update signature: %w", err)
	}

	log.Printf("Work paper signed with user: SignatureID=%s, UserID=%s", signature.ID, signature.UserID)

	return signature, nil
}

func (s *deskService) ListWorkPaperSignatures(ctx context.Context, req *ListWorkPaperSignaturesRequest) (*ListWorkPaperSignaturesResponse, error) {
	// Build query params
	queryParams := &pagination.QueryParams{
		Pagination: pagination.Pagination{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Filters: []pagination.Filter{},
		Sorts:   []pagination.Sort{},
	}

	// Add filters
	if req.UserID != "" {
		queryParams.Filters = append(queryParams.Filters, pagination.Filter{
			Field:    "user_id",
			Operator: "eq",
			Value:    req.UserID,
		})
	}

	if req.Status != "" {
		queryParams.Filters = append(queryParams.Filters, pagination.Filter{
			Field:    "status",
			Operator: "eq",
			Value:    req.Status,
		})
	}

	if req.WorkPaperID != "" {
		queryParams.Filters = append(queryParams.Filters, pagination.Filter{
			Field:    "work_paper_id",
			Operator: "eq",
			Value:    req.WorkPaperID,
		})
	}

	// Add sorting
	if req.SortBy != "" {
		order := "asc"
		if req.SortDirection == "desc" {
			order = "desc"
		}
		queryParams.Sorts = append(queryParams.Sorts, pagination.Sort{
			Field: req.SortBy,
			Order: order,
		})
	} else {
		// Default sort by created_at desc
		queryParams.Sorts = append(queryParams.Sorts, pagination.Sort{
			Field: "created_at",
			Order: "desc",
		})
	}

	// Get signatures
	signatures, total, err := s.signatureRepo.List(ctx, queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to list signatures: %w", err)
	}

	// Calculate pagination metadata
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &ListWorkPaperSignaturesResponse{
		Signatures:  signatures,
		TotalItems:  total,
		TotalPages:  totalPages,
		CurrentPage: req.Page,
		Limit:       req.Limit,
	}, nil
}

// Helper functions
