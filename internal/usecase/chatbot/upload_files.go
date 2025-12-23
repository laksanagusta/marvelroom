package chatbot

import (
	"context"
	"fmt"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/internal/infrastructure/gemini"
)

// UploadFilesUseCase handles uploading files to a knowledge base
type UploadFilesUseCase struct {
	repo             repository.ChatbotRepository
	fileSearchClient *gemini.FileSearchClient
}

// NewUploadFilesUseCase creates a new instance
func NewUploadFilesUseCase(repo repository.ChatbotRepository, fileSearchClient *gemini.FileSearchClient) *UploadFilesUseCase {
	return &UploadFilesUseCase{
		repo:             repo,
		fileSearchClient: fileSearchClient,
	}
}

// UploadFileInput represents a single file to upload
type UploadFileInput struct {
	FileName string
	Content  []byte
	MimeType string
}

// UploadFilesInput is the input for uploading files
type UploadFilesInput struct {
	KnowledgeBaseID string
	Files           []UploadFileInput
}

// UploadFilesOutput is the output
type UploadFilesOutput struct {
	Documents []*entity.KnowledgeBaseDocument
	Errors    []error
}

// Execute uploads files to the knowledge base
// Only users with "Super Admin" role can upload files
func (uc *UploadFilesUseCase) Execute(ctx context.Context, input UploadFilesInput, authenticatedUser entity.AuthenticatedUser) (*UploadFilesOutput, error) {
	// Check if user has super admin role
	if !authenticatedUser.HasRole("Super Admin") {
		return nil, fmt.Errorf("unauthorized: only super admin can upload files")
	}

	// Validate knowledge base exists
	kb, err := uc.repo.GetKnowledgeBaseByID(ctx, input.KnowledgeBaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get knowledge base: %w", err)
	}
	if kb == nil {
		return nil, fmt.Errorf("knowledge base not found")
	}

	output := &UploadFilesOutput{
		Documents: make([]*entity.KnowledgeBaseDocument, 0),
		Errors:    make([]error, 0),
	}

	var totalBytes int64 = 0

	for _, file := range input.Files {
		// Upload to Gemini FileSearchStore
		doc, err := uc.fileSearchClient.UploadToFileSearchStore(
			ctx,
			kb.FileSearchStoreID,
			file.FileName,
			file.Content,
			file.MimeType,
		)
		if err != nil {
			output.Errors = append(output.Errors, fmt.Errorf("failed to upload %s: %w", file.FileName, err))
			continue
		}

		// Create document record
		docEntity := &entity.KnowledgeBaseDocument{
			KnowledgeBaseID: input.KnowledgeBaseID,
			DocumentID:      doc.Name,
			FileName:        file.FileName,
			FileSize:        int64(len(file.Content)),
			MimeType:        file.MimeType,
			Status:          doc.State,
		}

		if err := uc.repo.CreateDocument(ctx, docEntity); err != nil {
			output.Errors = append(output.Errors, fmt.Errorf("failed to save document record for %s: %w", file.FileName, err))
			continue
		}

		output.Documents = append(output.Documents, docEntity)
		totalBytes += docEntity.FileSize
	}

	// Update knowledge base stats
	kb.TotalFiles += len(output.Documents)
	kb.TotalBytes += totalBytes
	if err := uc.repo.UpdateKnowledgeBase(ctx, kb); err != nil {
		// Non-critical error, just log it
		output.Errors = append(output.Errors, fmt.Errorf("failed to update knowledge base stats: %w", err))
	}

	return output, nil
}
