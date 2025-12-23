package chatbot

import (
	"context"
	"fmt"
	"strings"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/internal/infrastructure/gemini"
)

// SyncDocumentStatusUseCase handles syncing document status from Gemini
type SyncDocumentStatusUseCase struct {
	repo             repository.ChatbotRepository
	fileSearchClient *gemini.FileSearchClient
}

// NewSyncDocumentStatusUseCase creates a new instance
func NewSyncDocumentStatusUseCase(repo repository.ChatbotRepository, fileSearchClient *gemini.FileSearchClient) *SyncDocumentStatusUseCase {
	return &SyncDocumentStatusUseCase{
		repo:             repo,
		fileSearchClient: fileSearchClient,
	}
}

// SyncDocumentStatusInput is the input
type SyncDocumentStatusInput struct {
	KnowledgeBaseID string
}

// SyncDocumentStatusOutput is the output
type SyncDocumentStatusOutput struct {
	Documents      []*entity.KnowledgeBaseDocument
	UpdatedCount   int
	ProcessingDocs []string
	ActiveDocs     []string
	FailedDocs     []string
}

// Execute syncs document status from Gemini API to local database
// Only users with "Super Admin" role can sync document status
func (uc *SyncDocumentStatusUseCase) Execute(ctx context.Context, input SyncDocumentStatusInput, authenticatedUser entity.AuthenticatedUser) (*SyncDocumentStatusOutput, error) {
	// Check if user has super admin role
	if !authenticatedUser.HasRole("Super Admin") {
		return nil, fmt.Errorf("unauthorized: only super admin can sync document status")
	}

	// Get knowledge base
	kb, err := uc.repo.GetKnowledgeBaseByID(ctx, input.KnowledgeBaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get knowledge base: %w", err)
	}
	if kb == nil {
		return nil, fmt.Errorf("knowledge base not found")
	}

	// Get documents from our database
	localDocs, err := uc.repo.GetDocumentsByKnowledgeBaseID(ctx, input.KnowledgeBaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get local documents: %w", err)
	}

	// Get documents from Gemini
	geminiDocs, err := uc.fileSearchClient.ListDocuments(ctx, kb.FileSearchStoreID)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents from Gemini: %w", err)
	}

	// Create map of Gemini docs by name for quick lookup
	geminiDocMap := make(map[string]gemini.FileSearchDocument)
	for _, doc := range geminiDocs {
		geminiDocMap[doc.Name] = doc
	}

	output := &SyncDocumentStatusOutput{
		Documents:      localDocs,
		ProcessingDocs: []string{},
		ActiveDocs:     []string{},
		FailedDocs:     []string{},
	}

	// Update local documents with Gemini status
	for _, localDoc := range localDocs {
		// Find matching Gemini doc
		var newStatus string
		if geminiDoc, found := geminiDocMap[localDoc.DocumentID]; found {
			newStatus = strings.ToLower(geminiDoc.State)
		} else {
			// Document not found in Gemini, might have been deleted or name mismatch
			// Try to find by display name
			for _, gDoc := range geminiDocs {
				if gDoc.DisplayName == localDoc.FileName {
					newStatus = strings.ToLower(gDoc.State)
					// Update document ID to correct one
					localDoc.DocumentID = gDoc.Name
					break
				}
			}
			if newStatus == "" {
				newStatus = "failed" // Document not found in Gemini
			}
		}

		// Update if status changed
		if localDoc.Status != newStatus {
			if err := uc.repo.UpdateDocumentStatus(ctx, localDoc.ID, newStatus); err != nil {
				// Log but continue
				fmt.Printf("Warning: failed to update document status: %v\n", err)
			} else {
				output.UpdatedCount++
				localDoc.Status = newStatus
			}
		}

		// Categorize documents
		switch newStatus {
		case "processing":
			output.ProcessingDocs = append(output.ProcessingDocs, localDoc.FileName)
		case "active":
			output.ActiveDocs = append(output.ActiveDocs, localDoc.FileName)
		case "failed":
			output.FailedDocs = append(output.FailedDocs, localDoc.FileName)
		}
	}

	return output, nil
}
