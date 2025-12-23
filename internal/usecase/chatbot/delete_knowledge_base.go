package chatbot

import (
	"context"
	"fmt"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/internal/infrastructure/gemini"
)

// DeleteKnowledgeBaseUseCase handles deleting a knowledge base
type DeleteKnowledgeBaseUseCase struct {
	repo             repository.ChatbotRepository
	fileSearchClient *gemini.FileSearchClient
}

// NewDeleteKnowledgeBaseUseCase creates a new instance
func NewDeleteKnowledgeBaseUseCase(repo repository.ChatbotRepository, fileSearchClient *gemini.FileSearchClient) *DeleteKnowledgeBaseUseCase {
	return &DeleteKnowledgeBaseUseCase{
		repo:             repo,
		fileSearchClient: fileSearchClient,
	}
}

// DeleteKnowledgeBaseInput is the input
type DeleteKnowledgeBaseInput struct {
	KnowledgeBaseID string
}

// Execute deletes a knowledge base and its FileSearchStore
// Only users with "Super Admin" role can delete knowledge bases
func (uc *DeleteKnowledgeBaseUseCase) Execute(ctx context.Context, input DeleteKnowledgeBaseInput, authenticatedUser entity.AuthenticatedUser) error {
	// Check if user has super admin role
	if !authenticatedUser.HasRole("Super Admin") {
		return fmt.Errorf("unauthorized: only super admin can delete knowledge bases")
	}

	// Get knowledge base
	kb, err := uc.repo.GetKnowledgeBaseByID(ctx, input.KnowledgeBaseID)
	if err != nil {
		return fmt.Errorf("failed to get knowledge base: %w", err)
	}
	if kb == nil {
		return fmt.Errorf("knowledge base not found")
	}

	// First, delete all documents from FileSearchStore
	// This is required because Gemini doesn't allow deleting non-empty stores
	docs, err := uc.fileSearchClient.ListDocuments(ctx, kb.FileSearchStoreID)
	if err != nil {
		fmt.Printf("Warning: failed to list documents for deletion: %v\n", err)
	} else {
		for _, doc := range docs {
			if err := uc.fileSearchClient.DeleteDocument(ctx, kb.FileSearchStoreID, doc.Name); err != nil {
				fmt.Printf("Warning: failed to delete document %s: %v\n", doc.Name, err)
			}
		}
	}

	// Now delete the empty FileSearchStore from Gemini
	if err := uc.fileSearchClient.DeleteFileSearchStore(ctx, kb.FileSearchStoreID); err != nil {
		// Log but continue - we still want to clean up database
		fmt.Printf("Warning: failed to delete file search store: %v\n", err)
	}

	// Delete from database (cascades to documents and related chat sessions)
	if err := uc.repo.DeleteKnowledgeBase(ctx, input.KnowledgeBaseID); err != nil {
		return fmt.Errorf("failed to delete knowledge base: %w", err)
	}

	return nil
}
