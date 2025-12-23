package chatbot

import (
	"context"
	"fmt"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
)

// ListKnowledgeBasesUseCase handles listing user's knowledge bases
type ListKnowledgeBasesUseCase struct {
	repo repository.ChatbotRepository
}

// NewListKnowledgeBasesUseCase creates a new instance
func NewListKnowledgeBasesUseCase(repo repository.ChatbotRepository) *ListKnowledgeBasesUseCase {
	return &ListKnowledgeBasesUseCase{repo: repo}
}

// ListKnowledgeBasesOutput is the output
type ListKnowledgeBasesOutput struct {
	KnowledgeBases []*entity.KnowledgeBase
}

// Execute lists all global knowledge bases (shared across all users)
func (uc *ListKnowledgeBasesUseCase) Execute(ctx context.Context, authenticatedUser entity.AuthenticatedUser) (*ListKnowledgeBasesOutput, error) {
	// Return global knowledge bases (accessible by all users)
	kbs, err := uc.repo.GetGlobalKnowledgeBases(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list knowledge bases: %w", err)
	}

	return &ListKnowledgeBasesOutput{KnowledgeBases: kbs}, nil
}
