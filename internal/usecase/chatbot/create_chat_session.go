package chatbot

import (
	"context"
	"fmt"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
)

// CreateChatSessionUseCase handles creating a new chat session
type CreateChatSessionUseCase struct {
	repo repository.ChatbotRepository
}

// NewCreateChatSessionUseCase creates a new instance
func NewCreateChatSessionUseCase(repo repository.ChatbotRepository) *CreateChatSessionUseCase {
	return &CreateChatSessionUseCase{repo: repo}
}

// CreateChatSessionInput is the input
type CreateChatSessionInput struct {
	KnowledgeBaseID string
	Title           string
}

// CreateChatSessionOutput is the output
type CreateChatSessionOutput struct {
	ChatSession *entity.ChatSession
}

// Execute creates a new chat session
// Any authenticated user can create sessions on global knowledge bases
func (uc *CreateChatSessionUseCase) Execute(ctx context.Context, input CreateChatSessionInput, authenticatedUser entity.AuthenticatedUser) (*CreateChatSessionOutput, error) {
	// Validate knowledge base exists and is global
	kb, err := uc.repo.GetKnowledgeBaseByID(ctx, input.KnowledgeBaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get knowledge base: %w", err)
	}
	if kb == nil {
		return nil, fmt.Errorf("knowledge base not found")
	}
	// Only allow creating sessions on global knowledge bases
	if !kb.IsGlobal {
		return nil, fmt.Errorf("knowledge base not accessible")
	}

	title := input.Title
	if title == "" {
		title = "New Chat"
	}

	session := &entity.ChatSession{
		UserID:          authenticatedUser.ID,
		KnowledgeBaseID: input.KnowledgeBaseID,
		Title:           title,
	}

	if err := uc.repo.CreateChatSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create chat session: %w", err)
	}

	return &CreateChatSessionOutput{ChatSession: session}, nil
}
