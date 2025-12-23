package repository

import (
	"context"

	"sandbox/internal/domain/entity"
)

// ChatbotRepository defines the interface for chatbot data operations
type ChatbotRepository interface {
	// Knowledge Base operations
	CreateKnowledgeBase(ctx context.Context, kb *entity.KnowledgeBase) error
	GetKnowledgeBaseByID(ctx context.Context, id string) (*entity.KnowledgeBase, error)
	GetKnowledgeBasesByUserID(ctx context.Context, userID string) ([]*entity.KnowledgeBase, error)
	GetGlobalKnowledgeBases(ctx context.Context) ([]*entity.KnowledgeBase, error)
	UpdateKnowledgeBase(ctx context.Context, kb *entity.KnowledgeBase) error
	DeleteKnowledgeBase(ctx context.Context, id string) error

	// Knowledge Base Document operations
	CreateDocument(ctx context.Context, doc *entity.KnowledgeBaseDocument) error
	GetDocumentsByKnowledgeBaseID(ctx context.Context, kbID string) ([]*entity.KnowledgeBaseDocument, error)
	UpdateDocumentStatus(ctx context.Context, id string, status string) error
	DeleteDocument(ctx context.Context, id string) error

	// Chat Session operations
	CreateChatSession(ctx context.Context, session *entity.ChatSession) error
	GetChatSessionByID(ctx context.Context, id string) (*entity.ChatSession, error)
	GetChatSessionsByUserID(ctx context.Context, userID string) ([]*entity.ChatSession, error)
	UpdateChatSession(ctx context.Context, session *entity.ChatSession) error
	DeleteChatSession(ctx context.Context, id string) error

	// Chat Message operations
	CreateChatMessage(ctx context.Context, msg *entity.ChatMessage) error
	GetMessagesByChatSessionID(ctx context.Context, sessionID string) ([]*entity.ChatMessage, error)
}
