package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/pkg/database"

	"github.com/google/uuid"
)

type chatbotRepository struct {
	db database.DB
}

// NewChatbotRepository creates a new chatbot repository
func NewChatbotRepository(db database.DB) repository.ChatbotRepository {
	return &chatbotRepository{db: db}
}

// Knowledge Base operations

func (r *chatbotRepository) CreateKnowledgeBase(ctx context.Context, kb *entity.KnowledgeBase) error {
	if kb.ID == "" {
		kb.ID = uuid.New().String()
	}
	kb.CreatedAt = time.Now()
	kb.UpdatedAt = time.Now()

	query := `
		INSERT INTO knowledge_bases (id, user_id, name, file_search_store_id, total_files, total_bytes, is_global, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(ctx, query,
		kb.ID, kb.UserID, kb.Name, kb.FileSearchStoreID,
		kb.TotalFiles, kb.TotalBytes, kb.IsGlobal, kb.CreatedAt, kb.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create knowledge base: %w", err)
	}

	return nil
}

func (r *chatbotRepository) GetKnowledgeBaseByID(ctx context.Context, id string) (*entity.KnowledgeBase, error) {
	query := `
		SELECT id, user_id, name, file_search_store_id, total_files, total_bytes, is_global, created_at, updated_at
		FROM knowledge_bases
		WHERE id = $1`

	var kb entity.KnowledgeBase
	err := r.db.GetContext(ctx, &kb, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get knowledge base: %w", err)
	}

	return &kb, nil
}

func (r *chatbotRepository) GetKnowledgeBasesByUserID(ctx context.Context, userID string) ([]*entity.KnowledgeBase, error) {
	query := `
		SELECT id, user_id, name, file_search_store_id, total_files, total_bytes, is_global, created_at, updated_at
		FROM knowledge_bases
		WHERE user_id = $1
		ORDER BY created_at DESC`

	var kbs []*entity.KnowledgeBase
	err := r.db.SelectContext(ctx, &kbs, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list knowledge bases: %w", err)
	}

	return kbs, nil
}

func (r *chatbotRepository) GetGlobalKnowledgeBases(ctx context.Context) ([]*entity.KnowledgeBase, error) {
	query := `
		SELECT id, user_id, name, file_search_store_id, total_files, total_bytes, is_global, created_at, updated_at
		FROM knowledge_bases
		WHERE is_global = TRUE
		ORDER BY created_at DESC`

	var kbs []*entity.KnowledgeBase
	err := r.db.SelectContext(ctx, &kbs, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list global knowledge bases: %w", err)
	}

	return kbs, nil
}

func (r *chatbotRepository) UpdateKnowledgeBase(ctx context.Context, kb *entity.KnowledgeBase) error {
	kb.UpdatedAt = time.Now()

	query := `
		UPDATE knowledge_bases
		SET name = $1, total_files = $2, total_bytes = $3, updated_at = $4
		WHERE id = $5`

	_, err := r.db.ExecContext(ctx, query,
		kb.Name, kb.TotalFiles, kb.TotalBytes, kb.UpdatedAt, kb.ID)
	if err != nil {
		return fmt.Errorf("failed to update knowledge base: %w", err)
	}

	return nil
}

func (r *chatbotRepository) DeleteKnowledgeBase(ctx context.Context, id string) error {
	query := `DELETE FROM knowledge_bases WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete knowledge base: %w", err)
	}
	return nil
}

// Knowledge Base Document operations

func (r *chatbotRepository) CreateDocument(ctx context.Context, doc *entity.KnowledgeBaseDocument) error {
	if doc.ID == "" {
		doc.ID = uuid.New().String()
	}
	doc.CreatedAt = time.Now()

	query := `
		INSERT INTO knowledge_base_documents (id, knowledge_base_id, document_id, file_name, file_size, mime_type, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.ExecContext(ctx, query,
		doc.ID, doc.KnowledgeBaseID, doc.DocumentID,
		doc.FileName, doc.FileSize, doc.MimeType, doc.Status, doc.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}

	return nil
}

func (r *chatbotRepository) GetDocumentsByKnowledgeBaseID(ctx context.Context, kbID string) ([]*entity.KnowledgeBaseDocument, error) {
	query := `
		SELECT id, knowledge_base_id, document_id, file_name, file_size, mime_type, status, created_at
		FROM knowledge_base_documents
		WHERE knowledge_base_id = $1
		ORDER BY created_at DESC`

	var docs []*entity.KnowledgeBaseDocument
	err := r.db.SelectContext(ctx, &docs, query, kbID)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}

	return docs, nil
}

func (r *chatbotRepository) UpdateDocumentStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE knowledge_base_documents SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update document status: %w", err)
	}
	return nil
}

func (r *chatbotRepository) DeleteDocument(ctx context.Context, id string) error {
	query := `DELETE FROM knowledge_base_documents WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	return nil
}

// Chat Session operations

func (r *chatbotRepository) CreateChatSession(ctx context.Context, session *entity.ChatSession) error {
	if session.ID == "" {
		session.ID = uuid.New().String()
	}
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()

	query := `
		INSERT INTO chat_sessions (id, user_id, knowledge_base_id, title, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query,
		session.ID, session.UserID, session.KnowledgeBaseID,
		session.Title, session.CreatedAt, session.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create chat session: %w", err)
	}

	return nil
}

func (r *chatbotRepository) GetChatSessionByID(ctx context.Context, id string) (*entity.ChatSession, error) {
	query := `
		SELECT id, user_id, knowledge_base_id, title, created_at, updated_at
		FROM chat_sessions
		WHERE id = $1`

	var session entity.ChatSession
	err := r.db.GetContext(ctx, &session, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get chat session: %w", err)
	}

	return &session, nil
}

func (r *chatbotRepository) GetChatSessionsByUserID(ctx context.Context, userID string) ([]*entity.ChatSession, error) {
	query := `
		SELECT id, user_id, knowledge_base_id, title, created_at, updated_at
		FROM chat_sessions
		WHERE user_id = $1
		ORDER BY updated_at DESC`

	var sessions []*entity.ChatSession
	err := r.db.SelectContext(ctx, &sessions, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list chat sessions: %w", err)
	}

	return sessions, nil
}

func (r *chatbotRepository) UpdateChatSession(ctx context.Context, session *entity.ChatSession) error {
	session.UpdatedAt = time.Now()

	query := `UPDATE chat_sessions SET title = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, session.Title, session.UpdatedAt, session.ID)
	if err != nil {
		return fmt.Errorf("failed to update chat session: %w", err)
	}
	return nil
}

func (r *chatbotRepository) DeleteChatSession(ctx context.Context, id string) error {
	query := `DELETE FROM chat_sessions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete chat session: %w", err)
	}
	return nil
}

// Chat Message operations

// chatMessageDB is an internal struct for database scanning with JSON citations
type chatMessageDB struct {
	ID            string    `db:"id"`
	ChatSessionID string    `db:"chat_session_id"`
	Role          string    `db:"role"`
	Content       string    `db:"content"`
	Citations     string    `db:"citations"`
	CreatedAt     time.Time `db:"created_at"`
}

func (r *chatbotRepository) CreateChatMessage(ctx context.Context, msg *entity.ChatMessage) error {
	if msg.ID == "" {
		msg.ID = uuid.New().String()
	}
	msg.CreatedAt = time.Now()

	// Serialize citations to JSON
	citationsJSON := "[]"
	if len(msg.Citations) > 0 {
		data, err := json.Marshal(msg.Citations)
		if err != nil {
			return fmt.Errorf("failed to marshal citations: %w", err)
		}
		citationsJSON = string(data)
	}

	query := `
		INSERT INTO chat_messages (id, chat_session_id, role, content, citations, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query,
		msg.ID, msg.ChatSessionID, msg.Role, msg.Content, citationsJSON, msg.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create chat message: %w", err)
	}

	return nil
}

func (r *chatbotRepository) GetMessagesByChatSessionID(ctx context.Context, sessionID string) ([]*entity.ChatMessage, error) {
	query := `
		SELECT id, chat_session_id, role, content, citations, created_at
		FROM chat_messages
		WHERE chat_session_id = $1
		ORDER BY created_at ASC`

	var dbMessages []chatMessageDB
	err := r.db.SelectContext(ctx, &dbMessages, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}

	messages := make([]*entity.ChatMessage, 0, len(dbMessages))
	for _, dbMsg := range dbMessages {
		msg := &entity.ChatMessage{
			ID:            dbMsg.ID,
			ChatSessionID: dbMsg.ChatSessionID,
			Role:          dbMsg.Role,
			Content:       dbMsg.Content,
			CreatedAt:     dbMsg.CreatedAt,
		}

		// Parse citations JSON
		if dbMsg.Citations != "" && dbMsg.Citations != "[]" {
			if err := json.Unmarshal([]byte(dbMsg.Citations), &msg.Citations); err != nil {
				return nil, fmt.Errorf("failed to unmarshal citations: %w", err)
			}
		}

		messages = append(messages, msg)
	}

	return messages, nil
}
