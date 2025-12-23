package chatbot

import (
	"context"
	"fmt"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/internal/infrastructure/gemini"
)

// CreateKnowledgeBaseUseCase handles creating a new knowledge base
type CreateKnowledgeBaseUseCase struct {
	repo             repository.ChatbotRepository
	fileSearchClient *gemini.FileSearchClient
}

// NewCreateKnowledgeBaseUseCase creates a new instance
func NewCreateKnowledgeBaseUseCase(repo repository.ChatbotRepository, fileSearchClient *gemini.FileSearchClient) *CreateKnowledgeBaseUseCase {
	return &CreateKnowledgeBaseUseCase{
		repo:             repo,
		fileSearchClient: fileSearchClient,
	}
}

// CreateKnowledgeBaseInput is the input for creating a knowledge base
type CreateKnowledgeBaseInput struct {
	Name string
}

// CreateKnowledgeBaseOutput is the output
type CreateKnowledgeBaseOutput struct {
	KnowledgeBase *entity.KnowledgeBase
}

// Execute creates a new knowledge base with a Gemini FileSearchStore
// Only users with "Super Admin" role can create knowledge bases
func (uc *CreateKnowledgeBaseUseCase) Execute(ctx context.Context, input CreateKnowledgeBaseInput, authenticatedUser entity.AuthenticatedUser) (*CreateKnowledgeBaseOutput, error) {
	// Check if user has super admin role
	if !authenticatedUser.HasRole("Super Admin") {
		return nil, fmt.Errorf("unauthorized: only super admin can create knowledge bases")
	}

	// Create FileSearchStore in Gemini
	store, err := uc.fileSearchClient.CreateFileSearchStore(ctx, input.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create file search store: %w", err)
	}

	// Create knowledge base record - always global for super admin created KBs
	kb := &entity.KnowledgeBase{
		UserID:            authenticatedUser.ID,
		Name:              input.Name,
		FileSearchStoreID: store.Name,
		TotalFiles:        0,
		TotalBytes:        0,
		IsGlobal:          true, // All super admin created KBs are global
	}

	if err := uc.repo.CreateKnowledgeBase(ctx, kb); err != nil {
		// Cleanup: delete the FileSearchStore if database insert fails
		_ = uc.fileSearchClient.DeleteFileSearchStore(ctx, store.Name)
		return nil, fmt.Errorf("failed to create knowledge base: %w", err)
	}

	return &CreateKnowledgeBaseOutput{KnowledgeBase: kb}, nil
}
