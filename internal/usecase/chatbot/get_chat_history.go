package chatbot

import (
	"context"
	"fmt"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
)

// GetChatHistoryUseCase handles retrieving chat messages
type GetChatHistoryUseCase struct {
	repo repository.ChatbotRepository
}

// NewGetChatHistoryUseCase creates a new instance
func NewGetChatHistoryUseCase(repo repository.ChatbotRepository) *GetChatHistoryUseCase {
	return &GetChatHistoryUseCase{repo: repo}
}

// GetChatHistoryInput is the input
type GetChatHistoryInput struct {
	ChatSessionID string
}

// GetChatHistoryOutput is the output
type GetChatHistoryOutput struct {
	Messages []*entity.ChatMessage
}

// Execute retrieves chat messages for a session
func (uc *GetChatHistoryUseCase) Execute(ctx context.Context, input GetChatHistoryInput, authenticatedUser entity.AuthenticatedUser) (*GetChatHistoryOutput, error) {
	// Validate session exists and belongs to user
	session, err := uc.repo.GetChatSessionByID(ctx, input.ChatSessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat session: %w", err)
	}
	if session == nil {
		return nil, fmt.Errorf("chat session not found")
	}
	if session.UserID != authenticatedUser.ID {
		return nil, fmt.Errorf("unauthorized access to chat session")
	}

	messages, err := uc.repo.GetMessagesByChatSessionID(ctx, input.ChatSessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	return &GetChatHistoryOutput{Messages: messages}, nil
}
