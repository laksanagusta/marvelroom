package chatbot

import (
	"context"
	"fmt"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
)

// DeleteChatSessionUseCase handles deleting a chat session
type DeleteChatSessionUseCase struct {
	repo repository.ChatbotRepository
}

// NewDeleteChatSessionUseCase creates a new instance
func NewDeleteChatSessionUseCase(repo repository.ChatbotRepository) *DeleteChatSessionUseCase {
	return &DeleteChatSessionUseCase{repo: repo}
}

// DeleteChatSessionInput is the input
type DeleteChatSessionInput struct {
	SessionID string
}

// Execute deletes a chat session
// Only the session owner can delete their session
func (uc *DeleteChatSessionUseCase) Execute(ctx context.Context, input DeleteChatSessionInput, authenticatedUser entity.AuthenticatedUser) error {
	// Get the session to verify ownership
	session, err := uc.repo.GetChatSessionByID(ctx, input.SessionID)
	if err != nil {
		return fmt.Errorf("failed to get chat session: %w", err)
	}
	if session == nil {
		return fmt.Errorf("chat session not found")
	}

	// Verify the user owns this session
	if session.UserID != authenticatedUser.ID {
		return fmt.Errorf("unauthorized: you can only delete your own chat sessions")
	}

	// Delete the session (this will cascade to delete all messages)
	if err := uc.repo.DeleteChatSession(ctx, input.SessionID); err != nil {
		return fmt.Errorf("failed to delete chat session: %w", err)
	}

	return nil
}
