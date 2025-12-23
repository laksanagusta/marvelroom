package chatbot

import (
	"context"
	"fmt"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/internal/infrastructure/gemini"
)

// SendMessageUseCase handles sending a message and getting RAG response
type SendMessageUseCase struct {
	repo             repository.ChatbotRepository
	fileSearchClient *gemini.FileSearchClient
}

// NewSendMessageUseCase creates a new instance
func NewSendMessageUseCase(repo repository.ChatbotRepository, fileSearchClient *gemini.FileSearchClient) *SendMessageUseCase {
	return &SendMessageUseCase{
		repo:             repo,
		fileSearchClient: fileSearchClient,
	}
}

// SendMessageInput is the input for sending a message
type SendMessageInput struct {
	ChatSessionID string
	Content       string
}

// SendMessageOutput is the output
type SendMessageOutput struct {
	UserMessage      *entity.ChatMessage
	AssistantMessage *entity.ChatMessage
}

// Execute sends a message and gets a response with RAG
func (uc *SendMessageUseCase) Execute(ctx context.Context, input SendMessageInput, authenticatedUser entity.AuthenticatedUser) (*SendMessageOutput, error) {
	// Get chat session
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

	// Get knowledge base for FileSearchStore
	kb, err := uc.repo.GetKnowledgeBaseByID(ctx, session.KnowledgeBaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get knowledge base: %w", err)
	}
	if kb == nil {
		return nil, fmt.Errorf("knowledge base not found")
	}

	// Save user message first
	userMsg := &entity.ChatMessage{
		ChatSessionID: input.ChatSessionID,
		Role:          entity.ChatRoleUser,
		Content:       input.Content,
	}
	if err := uc.repo.CreateChatMessage(ctx, userMsg); err != nil {
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}

	// Get chat history for context
	messages, err := uc.repo.GetMessagesByChatSessionID(ctx, input.ChatSessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat history: %w", err)
	}

	// Convert to Gemini chat format
	chatHistory := make([]gemini.ChatMessage, 0, len(messages)-1) // Exclude the message we just added
	for _, msg := range messages {
		if msg.ID == userMsg.ID {
			continue // Skip the current message
		}
		chatHistory = append(chatHistory, gemini.ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Query Gemini with RAG
	response, err := uc.fileSearchClient.QueryWithFileSearch(
		ctx,
		kb.FileSearchStoreID,
		input.Content,
		chatHistory,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI response: %w", err)
	}

	// Convert citations
	citations := make([]entity.Citation, 0, len(response.Citations))
	for _, c := range response.Citations {
		citations = append(citations, entity.Citation{
			DocumentName: c.DocumentName,
			Content:      c.Content,
			StartIndex:   c.StartIndex,
			EndIndex:     c.EndIndex,
		})
	}

	// Save assistant message
	assistantMsg := &entity.ChatMessage{
		ChatSessionID: input.ChatSessionID,
		Role:          entity.ChatRoleAssistant,
		Content:       response.Text,
		Citations:     citations,
	}
	if err := uc.repo.CreateChatMessage(ctx, assistantMsg); err != nil {
		return nil, fmt.Errorf("failed to save assistant message: %w", err)
	}

	// Update session title if this is the first message
	if len(messages) <= 1 {
		// Use first few words of user message as title
		title := input.Content
		if len(title) > 50 {
			title = title[:50] + "..."
		}
		session.Title = title
		_ = uc.repo.UpdateChatSession(ctx, session)
	}

	return &SendMessageOutput{
		UserMessage:      userMsg,
		AssistantMessage: assistantMsg,
	}, nil
}
