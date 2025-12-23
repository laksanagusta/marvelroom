package entity

import (
	"time"
)

// KnowledgeBase represents a user's knowledge store for RAG
type KnowledgeBase struct {
	ID                string    `json:"id" db:"id"`
	UserID            string    `json:"user_id" db:"user_id"`
	Name              string    `json:"name" db:"name"`
	FileSearchStoreID string    `json:"file_search_store_id" db:"file_search_store_id"`
	TotalFiles        int       `json:"total_files" db:"total_files"`
	TotalBytes        int64     `json:"total_bytes" db:"total_bytes"`
	IsGlobal          bool      `json:"is_global" db:"is_global"` // If true, accessible by all users
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// KnowledgeBaseDocument represents a document in a knowledge base
type KnowledgeBaseDocument struct {
	ID              string    `json:"id" db:"id"`
	KnowledgeBaseID string    `json:"knowledge_base_id" db:"knowledge_base_id"`
	DocumentID      string    `json:"document_id" db:"document_id"` // Gemini document ID
	FileName        string    `json:"file_name" db:"file_name"`
	FileSize        int64     `json:"file_size" db:"file_size"`
	MimeType        string    `json:"mime_type" db:"mime_type"`
	Status          string    `json:"status" db:"status"` // "processing", "active", "failed"
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// ChatSession represents a conversation session
type ChatSession struct {
	ID              string    `json:"id" db:"id"`
	UserID          string    `json:"user_id" db:"user_id"`
	KnowledgeBaseID string    `json:"knowledge_base_id" db:"knowledge_base_id"`
	Title           string    `json:"title" db:"title"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// ChatMessage represents a single message in a chat session
type ChatMessage struct {
	ID            string     `json:"id" db:"id"`
	ChatSessionID string     `json:"chat_session_id" db:"chat_session_id"`
	Role          string     `json:"role" db:"role"` // "user" or "assistant"
	Content       string     `json:"content" db:"content"`
	Citations     []Citation `json:"citations,omitempty" db:"-"` // Stored as JSONB
	CitationsJSON string     `json:"-" db:"citations"`           // Raw JSONB for database
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

// Citation represents a source reference from RAG
type Citation struct {
	DocumentName string `json:"document_name"`
	Content      string `json:"content"`
	StartIndex   int    `json:"start_index,omitempty"`
	EndIndex     int    `json:"end_index,omitempty"`
}

// DocumentStatus constants
const (
	DocumentStatusProcessing = "processing"
	DocumentStatusActive     = "active"
	DocumentStatusFailed     = "failed"
)

// ChatRole constants
const (
	ChatRoleUser      = "user"
	ChatRoleAssistant = "assistant"
)
