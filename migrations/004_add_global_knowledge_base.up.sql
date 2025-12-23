-- Add is_global flag to knowledge_bases table
-- When true, the knowledge base is shared across all users (managed by super admin)
-- When false, it's a personal knowledge base (legacy behavior)

ALTER TABLE knowledge_bases ADD COLUMN is_global BOOLEAN DEFAULT FALSE;

-- Index for faster queries on global knowledge bases
CREATE INDEX idx_knowledge_bases_is_global ON knowledge_bases(is_global);
