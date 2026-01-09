package service

import (
	"context"

	"sandbox/internal/domain/entity"
	"sandbox/pkg/pagination"

	"github.com/google/uuid"
	"github.com/invopop/validation"
)

// DriveService defines the interface for Google Drive operations
type DriveService interface {
	GetFilesFromFolder(ctx context.Context, folderLink string) ([]*DriveFile, error)
	GetSubfolders(ctx context.Context, folderLink string) ([]*DriveFolder, error)
	GetFolderByName(ctx context.Context, parentFolderLink, folderName string) (*DriveFolder, error)
	DownloadFile(ctx context.Context, fileID string) ([]byte, error)
}

// DriveFile represents a file from Google Drive
type DriveFile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"` // "document", "spreadsheet", "pdf", etc.
	URL  string `json:"url"`
}

// DriveFolder represents a folder from Google Drive
type DriveFolder struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	FileCount int    `json:"file_count"` // Number of files in this folder
}

// LLMService defines the interface for LLM operations
type LLMService interface {
	CheckDocument(ctx context.Context, req *DocumentCheckRequest) (*DocumentCheckResponse, error)
}

// DocumentCheckRequest represents the request for document checking
type DocumentCheckRequest struct {
	Number           string         `json:"number"`
	Classification   string         `json:"classification"`
	TopicDescription string         `json:"topic_description"`
	DeskInstruction  string         `json:"desk_instruction"`
	Documents        []DocumentFile `json:"documents"`
}

// DocumentFile represents a document file for LLM processing
type DocumentFile struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
	Type string `json:"type"`
}

// DocumentCheckResponse represents the response from document checking
type DocumentCheckResponse struct {
	IsValid bool        `json:"isValid"`
	Notes   string      `json:"notes"`
	Model   string      `json:"model"`
	Usage   *TokenUsage `json:"usage,omitempty"`
}

// TokenUsage represents token usage information
type TokenUsage struct {
	PromptTokens     int `json:"promptTokens"`
	CompletionTokens int `json:"completionTokens"`
	TotalTokens      int `json:"totalTokens"`
}

// DeskService defines the interface for desk operations
type DeskService interface {
	// Work Paper Item operations
	CreateWorkPaperItem(ctx context.Context, req *CreateWorkPaperItemRequest) (*entity.WorkPaperItem, error)
	GetWorkPaperItem(ctx context.Context, id string) (*entity.WorkPaperItem, error)
	UpdateWorkPaperItem(ctx context.Context, req *UpdateWorkPaperItemRequest) (*entity.WorkPaperItem, error)
	DeleteWorkPaperItem(ctx context.Context, id uuid.UUID) (*entity.WorkPaperItem, error)
	DeactivateWorkPaperItem(ctx context.Context, id string) error
	ActivateWorkPaperItem(ctx context.Context, id string) error
	ListWorkPaperItems(ctx context.Context, params *entity.PaginationParams) (*entity.WorkPaperItemListResponse, error)
	GetActiveWorkPaperItems(ctx context.Context) ([]*entity.WorkPaperItem, error)

	// Organization operations
	GetOrganizations(ctx context.Context, page, limit int, sort string) (*entity.OrganizationListResponse, error)
	GetOrganization(ctx context.Context, id string) (*entity.Organization, error)

	// Work Paper Topic operations
	CreateWorkPaperTopic(ctx context.Context, req *CreateWorkPaperTopicRequest) (*entity.WorkPaperTopic, error)
	GetWorkPaperTopic(ctx context.Context, id string) (*entity.WorkPaperTopic, error)
	UpdateWorkPaperTopic(ctx context.Context, req *UpdateWorkPaperTopicRequest) (*entity.WorkPaperTopic, error)
	DeleteWorkPaperTopic(ctx context.Context, id string) error
	ListWorkPaperTopics(ctx context.Context, params *pagination.QueryParams) ([]*entity.WorkPaperTopic, int64, error)
	GetActiveWorkPaperTopics(ctx context.Context) ([]*entity.WorkPaperTopic, error)
	UploadWorkPaperTopicTemplate(ctx context.Context, req *UploadWorkPaperTopicTemplateRequest) (*entity.WorkPaperTopic, error)
	DeleteWorkPaperTopicTemplate(ctx context.Context, topicID string) error

	// Work Paper operations
	CreateWorkPaper(ctx context.Context, req *CreateWorkPaperRequest) (*entity.WorkPaper, error)
	GetWorkPaper(ctx context.Context, id string) (*entity.WorkPaper, error)
	GetWorkPaperByOrganizationYearSemester(ctx context.Context, organizationID string, year, semester int) (*entity.WorkPaper, error)
	UpdateWorkPaper(ctx context.Context, req *UpdateWorkPaperRequest) (*entity.WorkPaper, error)
	UpdateWorkPaperStatus(ctx context.Context, id string, status string) error
	ListWorkPapers(ctx context.Context, params *pagination.QueryParams) ([]*entity.WorkPaper, int64, error)
	ListWorkPapersByOrganization(ctx context.Context, organizationID string) ([]*entity.WorkPaper, error)
	SyncWorkPaperFolder(ctx context.Context, workPaperID string) (*SyncFolderResponse, error)

	// Work Paper Note operations
	GetWorkPaperNotes(ctx context.Context, workPaperID string) ([]*entity.WorkPaperNote, error)
	GetWorkPaperNoteByID(ctx context.Context, noteID string) (*entity.WorkPaperNote, error)
	UpdateWorkPaperNoteLink(ctx context.Context, noteID string, driveLink string) (*entity.WorkPaperNote, error)
	CheckDocument(ctx context.Context, noteID string) (*CheckDocumentResponse, error)
	UpdateWorkPaperNoteValidation(ctx context.Context, noteID string, isValid *bool, notes string) (*entity.WorkPaperNote, error)

	// Work Paper Signature operations
	CreateWorkPaperSignature(ctx context.Context, req *CreateWorkPaperSignatureRequest) (*entity.WorkPaperSignature, error)
	GetWorkPaperSignature(ctx context.Context, signatureID string) (*entity.WorkPaperSignature, error)
	GetWorkPaperSignatures(ctx context.Context, workPaperID string) ([]*entity.WorkPaperSignature, error)
	SignWorkPaper(ctx context.Context, signatureID string, req *SignWorkPaperRequest) (*entity.WorkPaperSignature, error)
	SignWorkPaperWithUser(ctx context.Context, signatureID string, userID string) (*entity.WorkPaperSignature, error)
	RejectWorkPaperSignature(ctx context.Context, signatureID string, req *RejectWorkPaperSignatureRequest) (*entity.WorkPaperSignature, error)
	ResetWorkPaperSignature(ctx context.Context, signatureID string) (*entity.WorkPaperSignature, error)
	GetWorkPaperSignaturesByUserID(ctx context.Context, userID string) ([]*entity.WorkPaperSignature, error)
	GetPendingSignaturesByUserID(ctx context.Context, userID string) ([]*entity.WorkPaperSignature, error)
	GetWorkPapersWithSignatures(ctx context.Context, page, limit int, status, organizationID string) ([]*WorkPaperWithSignatures, error)
	ListWorkPaperSignatures(ctx context.Context, req *ListWorkPaperSignaturesRequest) (*ListWorkPaperSignaturesResponse, error)

	// Work Paper Signer Management operations
	ManageSigners(ctx context.Context, req *ManageSignersRequest) (*ManageSignersResponse, error)

	// Backward compatibility methods (deprecated)
	CreateMasterLakipItem(ctx context.Context, req *CreateMasterLakipItemRequest) (*entity.WorkPaperItem, error)
	GetMasterLakipItem(ctx context.Context, id string) (*entity.WorkPaperItem, error)
	UpdateMasterLakipItem(ctx context.Context, id string, req *UpdateMasterLakipItemRequest) (*entity.WorkPaperItem, error)
	DeactivateMasterLakipItem(ctx context.Context, id string) error
	ActivateMasterLakipItem(ctx context.Context, id string) error
	ListMasterLakipItems(ctx context.Context, req *entity.PaginationParams) (*entity.WorkPaperItemListResponse, error)
	GetActiveMasterLakipItems(ctx context.Context) ([]*entity.WorkPaperItem, error)
	CreatePaperWork(ctx context.Context, req *CreatePaperWorkRequest) (*entity.WorkPaper, error)
	GetPaperWork(ctx context.Context, id string) (*entity.WorkPaper, error)
	GetPaperWorkByOrganizationYearSemester(ctx context.Context, organizationID string, year, semester int) (*entity.WorkPaper, error)
	UpdatePaperWorkStatus(ctx context.Context, id string, status string) error
	ListPaperWorks(ctx context.Context, req *ListPaperWorksRequest) ([]*entity.WorkPaper, int64, error)
	ListPaperWorksByOrganization(ctx context.Context, organizationID string) ([]*entity.WorkPaper, error)
	UpdatePaperWorkItemLink(ctx context.Context, itemID string, driveLink string) (*entity.WorkPaperNote, error)
	UpdatePaperWorkItemValidation(ctx context.Context, itemID string, isValid *bool, notes string) (*entity.WorkPaperNote, error)
}

// Request/Response DTOs
type CreateWorkPaperItemRequest struct {
	Type               string     `json:"type" validate:"required"`
	Number             string     `json:"number" validate:"required"`
	TopicID            *uuid.UUID `json:"topic_id"`
	DeskInstruction    string     `json:"desk_instruction" validate:"required"`
	ExpectedFolderName *string    `json:"expected_folder_name"`
	ParentID           *uuid.UUID `json:"parent_id"`
	Level              int        `json:"level"`
}

type UpdateWorkPaperItemRequest struct {
	ID                 uuid.UUID  `json:"id" validate:"required"`
	Type               string     `json:"type" validate:"required"`
	Number             string     `json:"number" validate:"required"`
	TopicID            *uuid.UUID `json:"topic_id"`
	DeskInstruction    string     `json:"desk_instruction" validate:"required"`
	ExpectedFolderName *string    `json:"expected_folder_name"`
	ParentID           *uuid.UUID `json:"parent_id"`
	Level              int        `json:"level"`
	Sequence           *int       `json:"sequence"`
	IsActive           *bool      `json:"is_active"`
}

type ListWorkPaperItemsRequest struct {
	Search   string `json:"search"`
	IsActive *bool  `json:"is_active"`
}

// Work Paper Topic DTOs
type CreateWorkPaperTopicRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type UpdateWorkPaperTopicRequest struct {
	ID          string `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

// UploadWorkPaperTopicTemplateRequest represents the request for uploading an Excel template
type UploadWorkPaperTopicTemplateRequest struct {
	TopicID     string `json:"topic_id" validate:"required"`
	FileName    string `json:"file_name" validate:"required"`
	FileContent []byte `json:"file_content" validate:"required"`
	ContentType string `json:"content_type"`
}

type CreateWorkPaperRequest struct {
	OrganizationID   string `json:"organization_id" validate:"required"`
	Name             string `json:"name"`
	Year             int    `json:"year" validate:"required,min=2000,max=2100"`
	Semester         int    `json:"semester" validate:"required,oneof=1 2"`
	TopicID          string `json:"topic_id"`
	SourceFolderLink string `json:"source_folder_link"` // Google Drive root folder link
}

// UpdateWorkPaperRequest represents the request for updating work paper
type UpdateWorkPaperRequest struct {
	ID               string  `json:"id" validate:"required"`
	Name             *string `json:"name,omitempty"`
	Status           *string `json:"status,omitempty" validate:"omitempty,oneof=draft ongoing ready_to_sign completed"`
	SourceFolderLink *string `json:"source_folder_link,omitempty"`
}

type ListWorkPapersRequest struct {
	OrganizationID string `json:"organization_id"`
	Year           *int   `json:"year"`
	Semester       *int   `json:"semester"`
	Status         string `json:"status"`
}

type CheckDocumentResponse struct {
	IsValid bool   `json:"isValid"`
	Notes   string `json:"notes"`
	Model   string `json:"model"`
}

// SyncFolderResponse represents the response from folder sync operation
type SyncFolderResponse struct {
	WorkPaperID     string           `json:"work_paper_id"`
	SyncedAt        string           `json:"synced_at"`
	TotalNotes      int              `json:"total_notes"`
	FoundCount      int              `json:"found_count"`
	MissingCount    int              `json:"missing_count"`
	LinkedCount     int              `json:"linked_count"`
	NotesSyncStatus []NoteSyncStatus `json:"notes_sync_status"`
}

// NoteSyncStatus represents the sync status of a single note
type NoteSyncStatus struct {
	NoteID             string `json:"note_id"`
	MasterItemNumber   string `json:"master_item_number"`
	ExpectedFolderName string `json:"expected_folder_name"`
	FileStatus         string `json:"file_status"`
	FilesInFolder      int    `json:"files_in_folder"`
	GDriveLink         string `json:"gdrive_link,omitempty"`
	Message            string `json:"message,omitempty"`
}

// Work Paper Signature DTOs
type CreateWorkPaperSignatureRequest struct {
	WorkPaperID   string                `json:"work_paper_id" validate:"required"`
	UserID        string                `json:"user_id" validate:"required"`
	UserName      string                `json:"user_name" validate:"required"`
	UserEmail     string                `json:"user_email"`
	UserRole      string                `json:"user_role"`
	SignatureType string                `json:"signature_type" validate:"required,oneof=digital manual approval"`
	SignatureData *entity.SignatureData `json:"signature_data"`
}

type SignWorkPaperRequest struct {
	SignatureData *entity.SignatureData `json:"signature_data"`
	Notes         string                `json:"notes"`
}

type RejectWorkPaperSignatureRequest struct {
	Notes string `json:"notes" validate:"required"`
}

type SignatureStatsResponse struct {
	Total    int `json:"total"`
	Pending  int `json:"pending"`
	Signed   int `json:"signed"`
	Rejected int `json:"rejected"`
}

// WorkPaperWithSignatures represents a work paper with its associated signatures
type WorkPaperWithSignatures struct {
	*entity.WorkPaper
	Signatures []*entity.WorkPaperSignature `json:"signatures"`
}

// ManageSignersRequest represents the request for managing signers
type ManageSignersRequest struct {
	WorkPaperID string             `json:"work_paper_id" validate:"required"`
	Action      string             `json:"action" validate:"required,oneof=add remove replace"`
	Signers     []CreateSignerData `json:"signers" validate:"required"`
}

// CreateSignerData represents signer data for management operations
type CreateSignerData struct {
	UserID        string `json:"user_id" validate:"required"`
	UserName      string `json:"user_name" validate:"required"`
	UserEmail     string `json:"user_email,omitempty"`
	UserRole      string `json:"user_role,omitempty"`
	SignatureType string `json:"signature_type" validate:"required,oneof=digital manual approval"`
}

// ManageSignersResponse represents the response for managing signers
type ManageSignersResponse struct {
	WorkPaperID string           `json:"work_paper_id"`
	Action      string           `json:"action"`
	Signers     []SignerResponse `json:"signers"`
	Message     string           `json:"message"`
}

// SignerResponse represents a signer in the response
type SignerResponse struct {
	SignatureID   string `json:"signature_id"`
	UserID        string `json:"user_id"`
	UserName      string `json:"user_name"`
	UserEmail     string `json:"user_email,omitempty"`
	UserRole      string `json:"user_role,omitempty"`
	SignatureType string `json:"signature_type"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

// ListWorkPaperSignaturesRequest represents request for listing work paper signatures with pagination and filtering
type ListWorkPaperSignaturesRequest struct {
	Page          int    `json:"page" validate:"min=1"`
	Limit         int    `json:"limit" validate:"min=1,max=100"`
	UserID        string `json:"user_id"`
	Status        string `json:"status"`
	WorkPaperID   string `json:"work_paper_id"`
	SortBy        string `json:"sort_by"`
	SortDirection string `json:"sort_direction" validate:"oneof=asc desc"`
}

// ListWorkPaperSignaturesResponse represents response for listing work paper signatures
type ListWorkPaperSignaturesResponse struct {
	Signatures  []*entity.WorkPaperSignature `json:"signatures"`
	TotalItems  int64                        `json:"total_items"`
	TotalPages  int                          `json:"total_pages"`
	CurrentPage int                          `json:"current_page"`
	Limit       int                          `json:"limit"`
}

// Validate methods for request structs

func (req *CreateWorkPaperSignatureRequest) Validate() error {
	if req.WorkPaperID == "" {
		return validation.NewError("work_paper_id", "Work paper ID is required")
	}
	if req.UserID == "" {
		return validation.NewError("user_id", "User ID is required")
	}
	if req.UserName == "" {
		return validation.NewError("user_name", "User name is required")
	}
	if req.SignatureType == "" {
		return validation.NewError("signature_type", "Signature type is required")
	}

	return validation.ValidateStruct(req)
}

func (req *SignWorkPaperRequest) Validate() error {
	// No required fields for this request
	return validation.ValidateStruct(req)
}

func (req *RejectWorkPaperSignatureRequest) Validate() error {
	if req.Notes == "" {
		return validation.NewError("notes", "Notes are required for rejection")
	}
	return validation.ValidateStruct(req)
}

func (req *ManageSignersRequest) Validate() error {
	if req.WorkPaperID == "" {
		return validation.NewError("work_paper_id", "Work paper ID is required")
	}
	if req.Action == "" {
		return validation.NewError("action", "Action is required")
	}
	if len(req.Signers) == 0 {
		return validation.NewError("signers", "At least one signer is required")
	}
	return validation.ValidateStruct(req)
}

func (req *ListWorkPaperSignaturesRequest) Validate() error {
	if req.Page < 1 {
		return validation.NewError("page", "Page must be at least 1")
	}
	if req.Limit < 1 || req.Limit > 100 {
		return validation.NewError("limit", "Limit must be between 1 and 100")
	}
	if req.SortDirection != "" && req.SortDirection != "asc" && req.SortDirection != "desc" {
		return validation.NewError("sort_direction", "Sort direction must be 'asc' or 'desc'")
	}
	return validation.ValidateStruct(req)
}

// Backward compatibility aliases (deprecated)
type (
	CreateMasterLakipItemRequest = CreateWorkPaperItemRequest
	UpdateMasterLakipItemRequest = UpdateWorkPaperItemRequest
	ListMasterLakipItemsRequest  = ListWorkPaperItemsRequest
	CreatePaperWorkRequest       = CreateWorkPaperRequest
	ListPaperWorksRequest        = ListWorkPapersRequest
)
