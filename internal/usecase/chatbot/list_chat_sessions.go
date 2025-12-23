package chatbot

import (
	"context"
	"fmt"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
)

// ListChatSessionsUseCase handles listing user's chat sessions
type ListChatSessionsUseCase struct {
	repo repository.ChatbotRepository
}

// NewListChatSessionsUseCase creates a new instance
func NewListChatSessionsUseCase(repo repository.ChatbotRepository) *ListChatSessionsUseCase {
	return &ListChatSessionsUseCase{repo: repo}
}

// ListChatSessionsOutput is the output
type ListChatSessionsOutput struct {
	Sessions []*entity.ChatSession
}

// Execute lists all chat sessions for a user
func (uc *ListChatSessionsUseCase) Execute(ctx context.Context, authenticatedUser entity.AuthenticatedUser) (*ListChatSessionsOutput, error) {
	sessions, err := uc.repo.GetChatSessionsByUserID(ctx, authenticatedUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list chat sessions: %w", err)
	}

	return &ListChatSessionsOutput{Sessions: sessions}, nil
}
