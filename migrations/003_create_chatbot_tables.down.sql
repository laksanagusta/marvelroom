-- Rollback chatbot tables

DROP INDEX IF EXISTS idx_chat_messages_created_at;
DROP INDEX IF EXISTS idx_chat_messages_session_id;
DROP INDEX IF EXISTS idx_chat_sessions_kb_id;
DROP INDEX IF EXISTS idx_chat_sessions_user_id;
DROP INDEX IF EXISTS idx_kb_documents_kb_id;
DROP INDEX IF EXISTS idx_knowledge_bases_user_id;

DROP TABLE IF EXISTS chat_messages;
DROP TABLE IF EXISTS chat_sessions;
DROP TABLE IF EXISTS knowledge_base_documents;
DROP TABLE IF EXISTS knowledge_bases;
