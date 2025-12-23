package chatbot

import (
	"context"
	"fmt"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
)

// GetKnowledgeBaseUseCase handles getting a single knowledge base with documents
type GetKnowledgeBaseUseCase struct {
	repo repository.ChatbotRepository
}

// NewGetKnowledgeBaseUseCase creates a new instance
func NewGetKnowledgeBaseUseCase(repo repository.ChatbotRepository) *GetKnowledgeBaseUseCase {
	return &GetKnowledgeBaseUseCase{repo: repo}
}

// GetKnowledgeBaseInput is the input
type GetKnowledgeBaseInput struct {
	KnowledgeBaseID string
}

// GetKnowledgeBaseOutput is the output
type GetKnowledgeBaseOutput struct {
	KnowledgeBase *entity.KnowledgeBase
	Documents     []*entity.KnowledgeBaseDocument
}

// Execute gets a knowledge base with its documents
// Any authenticated user can view global knowledge bases
func (uc *GetKnowledgeBaseUseCase) Execute(ctx context.Context, input GetKnowledgeBaseInput, authenticatedUser entity.AuthenticatedUser) (*GetKnowledgeBaseOutput, error) {
	kb, err := uc.repo.GetKnowledgeBaseByID(ctx, input.KnowledgeBaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get knowledge base: %w", err)
	}
	if kb == nil {
		return nil, fmt.Errorf("knowledge base not found")
	}
	// Only allow access to global knowledge bases
	if !kb.IsGlobal {
		return nil, fmt.Errorf("knowledge base not accessible")
	}

	docs, err := uc.repo.GetDocumentsByKnowledgeBaseID(ctx, input.KnowledgeBaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}

	return &GetKnowledgeBaseOutput{
		KnowledgeBase: kb,
		Documents:     docs,
	}, nil
}
